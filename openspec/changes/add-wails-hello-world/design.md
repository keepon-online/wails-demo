## Context

本项目需要实现一个基于 Wails v2 框架的 Windows 桌面应用，包含以下技术特点：
- 使用 Go 作为后端
- 使用 Web 技术（HTML/CSS/JS）作为前端
- 集成自动更新功能，支持从 GitHub Releases 获取更新

### 约束条件
- 目标平台：Windows
- 自动更新需要处理 Windows 上的文件锁定问题
- 需要考虑网络环境对更新检查的影响

## Goals / Non-Goals

### Goals
- 创建可运行的 Wails Hello World 应用
- 实现基于 GitHub Releases 的自动更新机制
- 提供清晰的版本显示和更新提示 UI
- 支持后台检查更新和用户主动检查

### Non-Goals
- 不实现 macOS/Linux 平台支持（首期仅 Windows）
- 不实现增量更新（使用完整二进制替换）
- 不实现多渠道更新（仅使用 GitHub Releases）
- 不实现强制更新机制

## Decisions

### 技术选型

1. **Wails v2** - 轻量级跨平台桌面应用框架
   - 相比 Electron 更轻量（不嵌入 Chromium）
   - 使用系统原生 WebView
   - Go 后端性能优秀

2. **go-selfupdate 或 goself** - 自动更新库
   - 支持 GitHub Releases 作为更新源
   - 处理二进制替换和版本比较
   - 社区验证可与 Wails 配合使用

3. **Vanilla JavaScript + CSS** - 前端技术栈
   - 保持简单，便于理解和维护
   - 不引入额外框架复杂度

### 替代方案考虑

| 方案 | 优点 | 缺点 | 选择原因 |
|------|------|------|---------|
| Electron | 生态成熟 | 包体积大（100MB+） | 不选：过于重量级 |
| Tauri | 类似 Wails | 需要 Rust 知识 | 不选：用户要求 Wails |
| Sidecar 进程更新 | 进程隔离 | 复杂度高 | 不选：Hello World 级别不需要 |
| 内置 go-selfupdate | 简单直接 | 需处理文件锁 | **选择**：适合简单场景 |

## Risks / Trade-offs

### 风险

1. **Windows 文件锁定**
   - 风险：更新时可执行文件被锁定导致失败
   - 缓解：使用延迟替换策略，重启后完成更新

2. **网络问题**
   - 风险：检查更新或下载失败
   - 缓解：实现重试机制和优雅降级

3. **版本格式不一致**
   - 风险：语义化版本解析错误
   - 缓解：强制使用 semver 格式 (vX.Y.Z)

## Migration Plan

这是新项目，无迁移需求。

## Open Questions

1. **更新检查频率**：启动时检查还是定时检查？
   - 建议：启动时检查 + 用户手动检查

2. **GitHub 仓库配置**：使用哪个 GitHub 仓库存放 releases？
   - 需要用户提供或后续配置

3. **版本命名规范**：是否强制使用 `v` 前缀？
   - 建议：使用 `v1.0.0` 格式
