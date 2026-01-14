package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// Config 应用配置结构
type Config struct {
	MinimizeToTray bool `json:"minimizeToTray"` // 关闭时最小化到托盘
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		MinimizeToTray: false, // 默认关闭时不隐藏到托盘
	}
}

// Store 配置存储管理器
type Store struct {
	config   *Config
	filePath string
	mu       sync.RWMutex
}

// NewStore 创建配置存储实例
func NewStore() (*Store, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	store := &Store{
		config:   DefaultConfig(),
		filePath: configPath,
	}

	// 尝试加载现有配置
	if err := store.Load(); err != nil {
		// 如果文件不存在，使用默认配置并保存
		if os.IsNotExist(err) {
			if err := store.Save(); err != nil {
				return nil, err
			}
		}
		// 其他错误忽略，使用默认配置
	}

	return store, nil
}

// getConfigPath 获取配置文件路径
// Windows: %APPDATA%/WailsDemo/config.json
func getConfigPath() (string, error) {
	// 获取用户配置目录
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	// 创建应用配置目录
	appConfigDir := filepath.Join(configDir, "WailsDemo")
	if err := os.MkdirAll(appConfigDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(appConfigDir, "config.json"), nil
}

// Load 从文件加载配置
func (s *Store) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return err
	}

	config := DefaultConfig()
	if err := json.Unmarshal(data, config); err != nil {
		return err
	}

	s.config = config
	return nil
}

// Save 保存配置到文件
func (s *Store) Save() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.MarshalIndent(s.config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filePath, data, 0644)
}

// Get 获取当前配置（只读副本）
func (s *Store) Get() Config {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return *s.config
}

// GetMinimizeToTray 获取最小化到托盘设置
func (s *Store) GetMinimizeToTray() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.config.MinimizeToTray
}

// SetMinimizeToTray 设置最小化到托盘选项并保存
func (s *Store) SetMinimizeToTray(value bool) error {
	s.mu.Lock()
	s.config.MinimizeToTray = value
	s.mu.Unlock()

	return s.Save()
}

// Update 更新配置并保存
func (s *Store) Update(cfg Config) error {
	s.mu.Lock()
	s.config = &cfg
	s.mu.Unlock()

	return s.Save()
}
