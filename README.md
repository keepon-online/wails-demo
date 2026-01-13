# Wails Auto-Update Demo

一个基于 Wails v2 的 Windows 桌面应用示例，展示如何集成 **自动更新** 功能。可作为其他 Wails 项目集成自动更新的参考。

## ✨ 功能特性

- 🚀 基于 [go-selfupdate](https://github.com/creativeprojects/go-selfupdate) 的自动更新
- 📦 支持便携版（直接替换 exe）和安装版（NSIS 安装包）
- 🔄 一键重启完成更新
- 🎨 现代化 UI 设计（玻璃拟态风格）

## 📁 项目结构

```
wails-demo/
├── main.go                          # 应用入口
├── app.go                           # 应用方法绑定
├── internal/
│   └── updater/
│       └── updater.go               # 🔑 自动更新核心模块
├── frontend/
│   ├── index.html
│   └── src/
│       ├── main.js                  # 前端逻辑
│       └── style.css                # 样式
├── build/
│   └── nsis/
│       └── installer.nsi            # NSIS 安装脚本
├── .github/
│   └── workflows/
│       └── release.yml              # CI/CD 自动发布
├── build.bat                        # 本地构建脚本
└── wails.json                       # Wails 配置
```

## 🚀 快速开始

### 环境要求

- Go 1.21+
- Node.js 18+
- [Wails CLI](https://wails.io/docs/gettingstarted/installation)
- (可选) [NSIS](https://nsis.sourceforge.io/) - 用于生成安装包

### 开发模式

```bash
wails dev
```

### 构建发布

```bash
# Windows
build.bat v1.0.0
```

输出文件：
- `dist/wails-demo_windows_amd64.zip` - 便携版
- `dist/wails-demo-setup-1.0.0.exe` - 安装版

## 🔧 集成自动更新到您的项目

### 步骤 1：添加依赖

```bash
go get github.com/creativeprojects/go-selfupdate
go get golang.org/x/sys/windows/registry
```

### 步骤 2：复制更新模块

将 `internal/updater/updater.go` 复制到您的项目，并修改包名和常量：

```go
// 修改 GitHub 仓库信息
updater.NewUpdater("your-username", "your-repo")
```

### 步骤 3：绑定方法

在 `app.go` 中添加：

```go
import "your-project/internal/updater"

type App struct {
    ctx     context.Context
    updater *updater.Updater
}

func NewApp() *App {
    return &App{
        updater: updater.NewUpdater("your-username", "your-repo"),
    }
}

// 检查更新
func (a *App) CheckForUpdate() (*updater.UpdateInfo, error) {
    return a.updater.CheckForUpdate(a.ctx)
}

// 应用更新
func (a *App) ApplyUpdate() (*updater.UpdateProgress, error) {
    return a.updater.DownloadAndApplyUpdate(a.ctx)
}

// 重启应用
func (a *App) RestartApp() error {
    return updater.RestartApp()
}
```

### 步骤 4：前端调用

```javascript
import { CheckForUpdate, ApplyUpdate, RestartApp } from '../wailsjs/go/main/App';

// 检查更新
const info = await CheckForUpdate();
if (info.available) {
    console.log(`发现新版本: ${info.latestVersion}`);
}

// 下载并应用更新
const result = await ApplyUpdate();
if (result.needRestart) {
    await RestartApp();
}
```

### 步骤 5：配置 GitHub Release

确保 Release 资源命名格式：
- `{app-name}_windows_amd64.zip` (便携版)
- `{app-name}-setup-{version}.exe` (安装版)

ZIP 包内的 exe 名称应与仓库名一致。

### 步骤 6：版本号注入

构建时通过 ldflags 注入版本号：

```bash
wails build -ldflags "-X 'your-project/internal/updater.Version=v1.0.0'"
```

## 🔄 更新流程

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  检查更新   │ ──▶ │  下载更新   │ ──▶ │  立即重启   │
└─────────────┘     └─────────────┘     └─────────────┘
       │                   │                   │
       ▼                   ▼                   ▼
  GitHub API          下载 ZIP/EXE       启动新进程
  查询最新版本        解压替换文件        退出当前进程
```

## 📋 更新模式

| 模式 | 检测方式 | 更新行为 |
|------|---------|---------|
| **便携版** | 无注册表项 | 下载 ZIP → 解压替换 exe → 重启 |
| **安装版** | 有注册表项 | 下载安装包 → 静默运行 → 自动重启 |

## 🛠 GitHub Actions

项目包含自动构建和发布工作流 (`.github/workflows/release.yml`)：

1. 推送版本标签触发构建
2. 自动构建 Windows 应用
3. 生成便携版 ZIP 和 NSIS 安装包
4. 创建 GitHub Release

```bash
# 发布新版本
git tag v1.0.1
git push origin v1.0.1
```

## 📝 License

MIT License

## 🔗 相关链接

- [Wails 官方文档](https://wails.io/)
- [go-selfupdate](https://github.com/creativeprojects/go-selfupdate)
- [NSIS](https://nsis.sourceforge.io/)
