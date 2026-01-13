# Change: 添加 Wails Hello World 演示应用（含 Windows 自动更新）

## Why

用户需要一个基于 Wails 框架的简单桌面应用演示项目，能够在 Windows 上实现自动更新功能。这将作为后续开发的基础模板。

## What Changes

- **新增** 基础的 Wails v2 项目结构
- **新增** 简单的 Hello World 前端界面
- **新增** Go 后端绑定示例
- **新增** Windows 自动更新功能（基于 GitHub Releases）
- **新增** 版本管理和更新检查机制

## Impact

- Affected specs: `wails-app`（新增）, `auto-update`（新增）
- Affected code:
  - 项目根目录（Wails 项目文件）
  - `frontend/` - 前端资源
  - `internal/` - 后端逻辑
  - `build/` - 构建配置
