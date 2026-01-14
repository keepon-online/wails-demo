package main

import (
	"context"
	"embed"
	"fmt"

	"wails-demo/internal/config"
	"wails-demo/internal/tray"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// 初始化配置存储
	configStore, err := config.NewStore()
	if err != nil {
		fmt.Println("警告: 无法初始化配置存储:", err)
		// 使用默认配置继续运行
	}

	// 创建托盘管理器
	trayManager := tray.NewManager(configStore)

	// 创建应用实例
	app := NewApp(configStore, trayManager)

	// 创建应用配置
	err = wails.Run(&options.App{
		Title:  "wails-demo",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		OnBeforeClose: func(ctx context.Context) bool {
			// 检查是否应该最小化到托盘
			if trayManager.ShouldMinimizeToTray() {
				// 隐藏窗口而不是退出
				runtime.WindowHide(ctx)
				return true // 返回 true 阻止关闭
			}
			return false // 返回 false 允许正常关闭
		},
		OnShutdown: app.shutdown,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
