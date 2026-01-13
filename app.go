package main

import (
	"context"
	"fmt"
	"wails-demo/internal/updater"
)

// App 应用程序结构体
type App struct {
	ctx     context.Context
	updater *updater.Updater
}

// NewApp 创建新的应用程序实例
func NewApp() *App {
	// TODO: 请修改为您的 GitHub 用户名和仓库名
	return &App{
		updater: updater.NewUpdater("keepon-online", "wails-demo"),
	}
}

// startup 应用启动时调用，保存上下文以便后续使用
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet 返回个性化问候语
func (a *App) Greet(name string) string {
	return fmt.Sprintf("你好 %s，欢迎使用 Wails 演示应用！", name)
}

// GetVersion 获取当前应用版本
func (a *App) GetVersion() string {
	return a.updater.GetCurrentVersion()
}

// CheckForUpdate 检查是否有可用更新
func (a *App) CheckForUpdate() (*updater.UpdateInfo, error) {
	return a.updater.CheckForUpdate(a.ctx)
}

// ApplyUpdate 下载并应用更新
func (a *App) ApplyUpdate() (*updater.UpdateProgress, error) {
	return a.updater.DownloadAndApplyUpdate(a.ctx)
}

// GetPlatformInfo 获取当前平台信息
func (a *App) GetPlatformInfo() string {
	return a.updater.GetPlatformInfo()
}

// RestartApp 重启应用以完成更新
func (a *App) RestartApp() error {
	return updater.RestartApp()
}
