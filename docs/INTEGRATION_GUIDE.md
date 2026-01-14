# Wails 自动更新集成指南

本文档详细说明如何将自动更新功能集成到您的 Wails 项目中。

## 目录

1. [前置条件](#前置条件)
2. [核心模块说明](#核心模块说明)
3. [集成步骤](#集成步骤)
4. [配置说明](#配置说明)
5. [常见问题](#常见问题)

---

## 前置条件

### 依赖库

```bash
go get github.com/creativeprojects/go-selfupdate
go get golang.org/x/sys/windows/registry  # Windows 注册表支持
```

### GitHub Release 要求

- 使用语义化版本标签：`v1.0.0`、`v1.2.3`
- 资源命名格式：`{app-name}_{os}_{arch}.zip`
- 示例：`myapp_windows_amd64.zip`

---

## 核心模块说明

### updater.go 关键结构

```go
// UpdateInfo 更新信息
type UpdateInfo struct {
    Available      bool   // 是否有可用更新
    CurrentVersion string // 当前版本
    LatestVersion  string // 最新版本
    ReleaseNotes   string // 发布说明
    ReleaseURL     string // 发布页面 URL
}

// UpdateProgress 更新进度
type UpdateProgress struct {
    Status      string  // checking, downloading, ready, error
    Progress    float64 // 进度百分比 (0-100)
    Message     string  // 状态消息
    NeedRestart bool    // 是否需要重启
}
```

### 更新模式

| 模式 | 触发条件 | 行为 |
|------|---------|------|
| ModePortable | 无注册表项 | 直接替换 exe |
| ModeInstaller | 有注册表项 | 下载运行安装包 |

---

## 集成步骤

### 1. 复制更新模块

```bash
mkdir -p internal/updater
cp /path/to/wails-demo/internal/updater/updater.go internal/updater/
```

### 2. 修改配置

编辑 `updater.go`，修改以下内容：

```go
// 修改版本变量名（如果需要）
var Version = "dev"

// 修改注册表路径（用于检测安装版）
registry.OpenKey(registry.LOCAL_MACHINE, `Software\YourAppName`, registry.READ)
```

### 3. 绑定方法到 App

```go
// app.go
type App struct {
    ctx     context.Context
    updater *updater.Updater
}

func NewApp() *App {
    return &App{
        updater: updater.NewUpdater("github-username", "repo-name"),
    }
}

func (a *App) CheckForUpdate() (*updater.UpdateInfo, error) {
    return a.updater.CheckForUpdate(a.ctx)
}

func (a *App) ApplyUpdate() (*updater.UpdateProgress, error) {
    return a.updater.DownloadAndApplyUpdate(a.ctx)
}

func (a *App) RestartApp() error {
    return updater.RestartApp()
}
```

### 4. 前端集成

```javascript
// 检查更新
async function checkUpdate() {
    const info = await CheckForUpdate();
    if (info.available) {
        showNotification(`发现新版本 ${info.latestVersion}`);
    }
}

// 下载并应用
async function applyUpdate() {
    const result = await ApplyUpdate();
    if (result.needRestart) {
        await RestartApp();
    }
}
```

### 5. 构建配置

```bash
# 注入版本号
wails build -ldflags "-X 'yourapp/internal/updater.Version=v1.0.0'"
```

---

## 配置说明

### GitHub Actions 工作流

```yaml
# .github/workflows/release.yml
on:
  push:
    tags:
      - 'v*'

jobs:
  build:
    runs-on: windows-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
      - run: go install github.com/wailsapp/wails/v2/cmd/wails@latest
      - run: |
          $version = "${{ github.ref_name }}"
          wails build -ldflags "-X 'yourapp/internal/updater.Version=$version'"
      - uses: softprops/action-gh-release@v2
        with:
          files: build/bin/*.zip
```

### NSIS 安装脚本要点

```nsis
; 写入注册表供自动更新检测
WriteRegStr HKLM "Software\YourAppName" "InstallPath" "$INSTDIR"
WriteRegStr HKLM "Software\YourAppName" "Version" "${PRODUCT_VERSION}"
```

---

## 常见问题

### Q: 更新检查返回 "未找到匹配的 Release 资源"

**原因**：GitHub Release 资源命名不符合格式要求

**解决**：确保命名为 `{app-name}_{os}_{arch}.zip`，如 `myapp_windows_amd64.zip`

### Q: "invalid semantic version" 错误

**原因**：版本号格式不正确

**解决**：使用语义化版本格式 `v1.0.0`，或在代码中处理 `dev` 版本：

```go
if currentVersion == "dev" {
    currentVersion = "0.0.0"
}
```

### Q: 更新失败，文件被锁定

**原因**：Windows 不允许替换正在运行的 exe

**解决**：`go-selfupdate` 会自动处理，通过重命名旧文件再写入新文件

### Q: 安装版更新后未自动重启

**原因**：NSIS 静默安装需要配置

**解决**：安装脚本中启用静默模式并配置自启动

---

## 系统托盘集成

### 概述

系统托盘功能允许应用在关闭窗口时最小化到托盘而非退出，用户可通过托盘图标右键菜单操作应用。

### 依赖

```bash
go get github.com/getlantern/systray
```

### 1. 创建配置模块

```go
// internal/config/config.go
type Config struct {
    MinimizeToTray bool `json:"minimizeToTray"`
}

type Store struct {
    config   *Config
    filePath string
    mu       sync.RWMutex
}

func NewStore() (*Store, error) {
    // 配置存储在 %APPDATA%/YourApp/config.json
    configPath, err := getConfigPath()
    // ...
}
```

### 2. 创建托盘模块

```go
// internal/tray/tray.go
type Manager struct {
    configStore    *config.Store
    showWindowFunc func()
    hideWindowFunc func()
    quitFunc       func()
}

func (m *Manager) Run() {
    go systray.Run(m.onReady, m.onExit)
}
```

托盘图标需要嵌入到模块中：

```go
//go:embed icon.ico
var iconFS embed.FS
```

### 3. 集成到应用

```go
// main.go
func main() {
    configStore, _ := config.NewStore()
    trayManager := tray.NewManager(configStore)
    app := NewApp(configStore, trayManager)

    wails.Run(&options.App{
        OnStartup: app.startup,
        OnBeforeClose: func(ctx context.Context) bool {
            if trayManager.ShouldMinimizeToTray() {
                runtime.WindowHide(ctx)
                return true // 阻止关闭
            }
            return false
        },
        OnShutdown: app.shutdown,
    })
}
```

### 4. 前端设置界面

```javascript
import { GetTraySettings, SetTraySettings } from '../wailsjs/go/main/App';

// 获取设置
const settings = await GetTraySettings();
checkbox.checked = settings.minimizeToTray;

// 保存设置
await SetTraySettings({ minimizeToTray: checkbox.checked });
```

### 5. 托盘菜单设计

```
[应用图标]
├── 显示主窗口
├── ─────────
├── ✓ 关闭时最小化到托盘
├── ─────────
└── 退出
```

---

## 参考链接

- [go-selfupdate 文档](https://github.com/creativeprojects/go-selfupdate)
- [systray 库](https://github.com/getlantern/systray)
- [NSIS 手册](https://nsis.sourceforge.io/Docs/)
- [GitHub Actions](https://docs.github.com/actions)
