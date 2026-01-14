package tray

import (
	"embed"
	"wails-demo/internal/config"

	"github.com/getlantern/systray"
)

//go:embed icon.ico
var iconFS embed.FS

// Manager 托盘管理器
type Manager struct {
	configStore     *config.Store
	showWindowFunc  func()
	hideWindowFunc  func()
	quitFunc        func()
	menuMinimize    *systray.MenuItem
	onReadyCallback func()
}

// NewManager 创建托盘管理器
func NewManager(configStore *config.Store) *Manager {
	return &Manager{
		configStore: configStore,
	}
}

// SetWindowFuncs 设置窗口控制函数
func (m *Manager) SetWindowFuncs(show, hide, quit func()) {
	m.showWindowFunc = show
	m.hideWindowFunc = hide
	m.quitFunc = quit
}

// SetOnReady 设置托盘就绪回调
func (m *Manager) SetOnReady(callback func()) {
	m.onReadyCallback = callback
}

// Run 启动托盘（在独立 goroutine 中运行）
func (m *Manager) Run() {
	go systray.Run(m.onReady, m.onExit)
}

// onReady 托盘就绪时调用
func (m *Manager) onReady() {
	// 设置托盘图标
	iconData, err := iconFS.ReadFile("icon.ico")
	if err == nil {
		systray.SetIcon(iconData)
	}

	// 设置托盘提示文本
	systray.SetTitle("Wails Demo")
	systray.SetTooltip("Wails Demo - 点击显示主窗口")

	// 创建菜单项
	menuShow := systray.AddMenuItem("显示主窗口", "显示应用主窗口")
	systray.AddSeparator()

	// 最小化到托盘开关
	m.menuMinimize = systray.AddMenuItemCheckbox(
		"关闭时最小化到托盘",
		"启用后关闭窗口将最小化到托盘",
		m.configStore.GetMinimizeToTray(),
	)

	systray.AddSeparator()
	menuQuit := systray.AddMenuItem("退出", "完全退出应用")

	// 触发就绪回调
	if m.onReadyCallback != nil {
		m.onReadyCallback()
	}

	// 菜单事件处理
	go func() {
		for {
			select {
			case <-menuShow.ClickedCh:
				if m.showWindowFunc != nil {
					m.showWindowFunc()
				}
			case <-m.menuMinimize.ClickedCh:
				// 切换设置
				enabled := !m.configStore.GetMinimizeToTray()
				m.configStore.SetMinimizeToTray(enabled)
				if enabled {
					m.menuMinimize.Check()
				} else {
					m.menuMinimize.Uncheck()
				}
			case <-menuQuit.ClickedCh:
				if m.quitFunc != nil {
					m.quitFunc()
				}
				systray.Quit()
			}
		}
	}()
}

// onExit 托盘退出时调用
func (m *Manager) onExit() {
	// 清理资源
}

// ShouldMinimizeToTray 检查是否应该最小化到托盘
func (m *Manager) ShouldMinimizeToTray() bool {
	return m.configStore.GetMinimizeToTray()
}

// Quit 退出托盘
func (m *Manager) Quit() {
	systray.Quit()
}

// OnDoubleClick 处理双击事件（由 Windows 消息循环调用）
func (m *Manager) OnDoubleClick() {
	if m.showWindowFunc != nil {
		m.showWindowFunc()
	}
}

// SetMinimizeToTray 设置最小化到托盘选项
func (m *Manager) SetMinimizeToTray(enabled bool) error {
	err := m.configStore.SetMinimizeToTray(enabled)
	if err != nil {
		return err
	}

	// 更新菜单状态
	if m.menuMinimize != nil {
		if enabled {
			m.menuMinimize.Check()
		} else {
			m.menuMinimize.Uncheck()
		}
	}
	return nil
}

// GetMinimizeToTray 获取最小化到托盘设置
func (m *Manager) GetMinimizeToTray() bool {
	return m.configStore.GetMinimizeToTray()
}
