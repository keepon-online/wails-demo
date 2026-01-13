## ADDED Requirements

### Requirement: Wails 应用基础结构
系统 SHALL 提供一个基于 Wails v2 框架的完整项目结构，包含前端、后端和构建配置。

#### Scenario: 项目初始化成功
- **WHEN** 用户克隆或下载项目后
- **THEN** 可以使用 `wails dev` 命令启动开发模式
- **AND** 可以使用 `wails build` 命令构建 Windows 可执行文件

#### Scenario: 开发模式热重载
- **WHEN** 开发者修改前端代码
- **THEN** 浏览器自动刷新显示更改
- **AND** 后端代码修改后自动重新编译

---

### Requirement: Hello World 用户界面
系统 SHALL 提供一个简单直观的 Hello World 界面，展示 Wails 的前后端通信能力。

#### Scenario: 显示欢迎信息
- **WHEN** 应用启动完成
- **THEN** 显示 "Hello World" 标题
- **AND** 显示当前应用版本号

#### Scenario: 前后端交互演示
- **WHEN** 用户输入姓名并点击"打招呼"按钮
- **THEN** 调用 Go 后端方法
- **AND** 返回个性化问候语并显示在界面上

---

### Requirement: Go 后端绑定
系统 SHALL 提供 Go 后端方法绑定，允许前端 JavaScript 调用 Go 函数。

#### Scenario: 暴露 Greet 方法
- **WHEN** 前端调用 `window.go.main.App.Greet(name)` 
- **THEN** Go 后端执行 `App.Greet` 方法
- **AND** 返回格式化的问候字符串

#### Scenario: 暴露版本信息方法
- **WHEN** 前端调用 `window.go.main.App.GetVersion()`
- **THEN** Go 后端返回当前应用版本号
