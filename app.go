package main

import (
	"context"
	"fmt"

	"wails-demo/internal/config"
	"wails-demo/internal/tray"
	"wails-demo/internal/updater"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App 应用程序结构体
type App struct {
	ctx         context.Context
	updater     *updater.Updater
	configStore *config.Store
	trayManager *tray.Manager
}

// NewApp 创建新的应用程序实例
func NewApp(configStore *config.Store, trayManager *tray.Manager) *App {
	return &App{
		updater:     updater.NewUpdater("keepon-online", "wails-demo"),
		configStore: configStore,
		trayManager: trayManager,
	}
}

// startup 应用启动时调用，保存上下文以便后续使用
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// 设置托盘窗口控制函数
	a.trayManager.SetWindowFuncs(
		func() { runtime.WindowShow(a.ctx) }, // 显示窗口
		func() { runtime.WindowHide(a.ctx) }, // 隐藏窗口
		func() { runtime.Quit(a.ctx) },       // 退出应用
	)

	// 启动托盘
	a.trayManager.Run()
}

// shutdown 应用关闭时调用
func (a *App) shutdown(ctx context.Context) {
	// 清理托盘
	a.trayManager.Quit()
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

// ===== 托盘设置相关 API =====

// TraySettings 托盘设置结构（用于前端）
type TraySettings struct {
	MinimizeToTray bool `json:"minimizeToTray"`
}

// GetTraySettings 获取托盘设置
func (a *App) GetTraySettings() TraySettings {
	return TraySettings{
		MinimizeToTray: a.trayManager.GetMinimizeToTray(),
	}
}

// SetTraySettings 更新托盘设置
func (a *App) SetTraySettings(settings TraySettings) error {
	return a.trayManager.SetMinimizeToTray(settings.MinimizeToTray)
}
