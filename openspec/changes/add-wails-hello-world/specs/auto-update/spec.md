## ADDED Requirements

### Requirement: 自动更新检查
系统 SHALL 提供基于 GitHub Releases 的自动更新检查功能。

#### Scenario: 启动时检查更新
- **WHEN** 应用启动完成
- **THEN** 后台检查 GitHub Releases 获取最新版本
- **AND** 如果有新版本，显示更新提示通知

#### Scenario: 用户手动检查更新
- **WHEN** 用户点击"检查更新"按钮
- **THEN** 显示检查中状态
- **AND** 查询 GitHub Releases API 获取最新版本信息
- **AND** 显示检查结果（已是最新版 或 发现新版本）

#### Scenario: 网络错误处理
- **WHEN** 检查更新时网络不可用
- **THEN** 显示友好的错误提示
- **AND** 不阻塞应用正常使用

---

### Requirement: 更新下载和安装
系统 SHALL 提供更新下载和安装功能，支持替换当前可执行文件。

#### Scenario: 下载更新
- **WHEN** 用户确认下载新版本
- **THEN** 从 GitHub Releases 下载对应平台的资源文件
- **AND** 显示下载进度百分比

#### Scenario: 应用更新
- **WHEN** 下载完成且用户确认安装
- **THEN** 替换当前可执行文件为新版本
- **AND** 提示用户重启应用以完成更新

#### Scenario: 更新失败回滚
- **WHEN** 更新过程中发生错误
- **THEN** 保留原有可执行文件不变
- **AND** 显示错误信息和重试选项

---

### Requirement: 版本管理
系统 SHALL 提供语义化版本管理，支持版本比较和显示。

#### Scenario: 版本号格式
- **WHEN** 应用版本号被设置
- **THEN** 使用语义化版本格式（vX.Y.Z）
- **AND** 支持通过构建参数注入版本号

#### Scenario: 版本比较
- **WHEN** 检查到远程版本
- **THEN** 正确比较本地和远程版本大小
- **AND** 只有远程版本更高时才提示更新

---

### Requirement: 更新状态 UI
系统 SHALL 在界面上清晰展示更新相关状态。

#### Scenario: 显示当前版本
- **WHEN** 应用运行时
- **THEN** 在界面底部或设置中显示"当前版本：vX.Y.Z"

#### Scenario: 显示更新进度
- **WHEN** 正在下载更新
- **THEN** 显示进度条和已下载百分比
- **AND** 提供取消下载的选项

#### Scenario: 更新完成提示
- **WHEN** 更新包下载并准备完成
- **THEN** 显示"更新已就绪，重启应用以完成安装"提示
- **AND** 提供"立即重启"按钮
