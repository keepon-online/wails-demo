# System Tray Spec

## Overview

系统托盘功能规范，定义关闭时隐藏到托盘的可选行为。

---

## ADDED Requirements

### REQ-TRAY-001: 托盘图标显示

应用运行时应在系统托盘显示图标。

**Acceptance Criteria:**
- 托盘图标使用应用图标
- 图标悬停显示应用名称

#### Scenario: 应用启动显示托盘图标

```
Given 应用启动
When 应用加载完成
Then 系统托盘显示应用图标
```

---

### REQ-TRAY-002: 托盘菜单

托盘图标右键点击应显示菜单。

**Acceptance Criteria:**
- 菜单包含"显示主窗口"选项
- 菜单包含"关闭时最小化到托盘"开关
- 菜单包含"退出"选项

#### Scenario: 右键显示托盘菜单

```
Given 应用正在运行
When 用户右键点击托盘图标
Then 显示托盘菜单
And 菜单包含"显示主窗口"、"最小化到托盘开关"、"退出"
```

---

### REQ-TRAY-003: 关闭时隐藏到托盘

启用该选项后，关闭窗口不退出应用，而是隐藏到托盘。

**Acceptance Criteria:**
- 该功能默认关闭
- 用户可通过托盘菜单或设置界面开启
- 开启后点击关闭按钮隐藏窗口而非退出

#### Scenario: 启用最小化到托盘后关闭窗口

```
Given "最小化到托盘"选项已开启
When 用户点击窗口关闭按钮
Then 窗口隐藏
And 托盘图标保持显示
And 应用继续在后台运行
```

#### Scenario: 禁用最小化到托盘后关闭窗口

```
Given "最小化到托盘"选项已关闭
When 用户点击窗口关闭按钮
Then 应用正常退出
```

---

### REQ-TRAY-004: 双击托盘显示窗口

双击托盘图标应显示主窗口。

**Acceptance Criteria:**
- 窗口隐藏时，双击托盘图标显示窗口
- 窗口已显示时，双击托盘图标将窗口置顶

#### Scenario: 双击托盘图标恢复窗口

```
Given 窗口已隐藏到托盘
When 用户双击托盘图标
Then 主窗口显示并激活
```

---

### REQ-TRAY-005: 设置持久化

托盘相关设置应持久化存储。

**Acceptance Criteria:**
- 设置保存到本地配置文件
- 应用重启后设置保持

#### Scenario: 设置在重启后保持

```
Given 用户开启"最小化到托盘"选项
When 应用退出并重新启动
Then "最小化到托盘"选项仍为开启状态
```

---

## Dependencies

- `wails-app` spec (主应用窗口)

## Related Specs

无
