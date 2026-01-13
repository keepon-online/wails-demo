# Project Context

## Purpose

Wails Demo 是一个基于 Wails v2 框架的 Windows 桌面应用演示项目，展示如何构建一个简单的 Hello World 应用并集成自动更新功能。

## Tech Stack

- **后端**: Go 1.21+
- **前端**: Vanilla JavaScript + CSS
- **桌面框架**: Wails v2
- **自动更新**: go-selfupdate / goself

## Project Conventions

### Code Style

- **Go**: 遵循 Effective Go 和 Uber Go Style Guide
- **JavaScript**: ES6+ 语法，使用 const/let
- **CSS**: BEM 命名规范
- **注释**: 使用简体中文

### Architecture Patterns

- 前后端分离，通过 Wails 绑定通信
- 单一职责：每个 Go 结构体负责特定功能
- 配置外置：版本号和仓库地址通过构建参数注入

### Testing Strategy

- **开发模式测试**: `wails dev` 运行开发服务器
- **构建测试**: `wails build` 生成可执行文件
- **更新测试**: 手动创建测试 Release 验证更新流程

### Git Workflow

- 主分支: `main`
- 功能分支: `feature/*`
- 版本标签: `v1.0.0` 格式
- 提交规范: Conventional Commits

## Domain Context

- 桌面应用需要处理 Windows 文件系统权限
- 自动更新涉及替换正在运行的可执行文件
- GitHub Releases API 有访问频率限制

## Important Constraints

- 首期仅支持 Windows 平台
- 自动更新依赖 GitHub Releases 可访问性
- 版本号必须遵循语义化版本规范

## External Dependencies

- GitHub API - 用于检查和下载更新
- Wails CLI - 用于开发和构建
