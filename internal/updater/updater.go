// Package updater 提供基于 GitHub Releases 的自动更新功能
package updater

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/creativeprojects/go-selfupdate"
	"golang.org/x/sys/windows/registry"
)

// 版本信息（通过 ldflags 注入）
var (
	Version = "dev"
)

// UpdateMode 更新模式
type UpdateMode string

const (
	ModePortable  UpdateMode = "portable"  // 便携版：直接替换 exe
	ModeInstaller UpdateMode = "installer" // 安装版：下载并运行新安装包
)

// UpdateInfo 更新信息结构
type UpdateInfo struct {
	Available      bool   `json:"available"`      // 是否有可用更新
	CurrentVersion string `json:"currentVersion"` // 当前版本
	LatestVersion  string `json:"latestVersion"`  // 最新版本
	ReleaseNotes   string `json:"releaseNotes"`   // 发布说明
	ReleaseURL     string `json:"releaseUrl"`     // 发布页面 URL
	InstallerURL   string `json:"installerUrl"`   // 安装包下载 URL
	DebugInfo      string `json:"debugInfo"`      // 调试信息
}

// UpdateProgress 更新进度信息
type UpdateProgress struct {
	Status      string  `json:"status"`      // 状态：checking, downloading, ready, error
	Progress    float64 `json:"progress"`    // 进度百分比 (0-100)
	Message     string  `json:"message"`     // 状态消息
	NeedRestart bool    `json:"needRestart"` // 是否需要重启
}

// Updater 自动更新器
type Updater struct {
	owner   string // GitHub 用户名/组织名
	repo    string // 仓库名
	mode    UpdateMode
	source  selfupdate.Source
	latest  *selfupdate.Release
	tempDir string // 临时下载目录
}

// NewUpdater 创建新的更新器实例
// owner: GitHub 用户名或组织名
// repo: 仓库名
func NewUpdater(owner, repo string) *Updater {
	// 使用 GitHub 作为更新源
	source, _ := selfupdate.NewGitHubSource(selfupdate.GitHubConfig{})

	// 检测更新模式：如果有注册表项则为安装版
	mode := ModePortable
	if isInstalledVersion() {
		mode = ModeInstaller
	}

	return &Updater{
		owner:  owner,
		repo:   repo,
		mode:   mode,
		source: source,
	}
}

// isInstalledVersion 检查是否为安装版
func isInstalledVersion() bool {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `Software\Wails Demo`, registry.READ)
	if err != nil {
		return false
	}
	defer key.Close()

	installPath, _, err := key.GetStringValue("InstallPath")
	if err != nil || installPath == "" {
		return false
	}
	return true
}

// getInstallPath 获取安装路径
func getInstallPath() (string, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `Software\Wails Demo`, registry.READ)
	if err != nil {
		return "", fmt.Errorf("无法读取注册表: %w", err)
	}
	defer key.Close()

	installPath, _, err := key.GetStringValue("InstallPath")
	if err != nil {
		return "", fmt.Errorf("无法获取安装路径: %w", err)
	}
	return installPath, nil
}

// GetCurrentVersion 获取当前版本
func (u *Updater) GetCurrentVersion() string {
	return Version
}

// GetUpdateMode 获取更新模式
func (u *Updater) GetUpdateMode() string {
	return string(u.mode)
}

// CheckForUpdate 检查是否有可用更新
func (u *Updater) CheckForUpdate(ctx context.Context) (*UpdateInfo, error) {
	info := &UpdateInfo{
		Available:      false,
		CurrentVersion: Version,
	}

	// 移除版本号前缀 v
	currentVersion := strings.TrimPrefix(Version, "v")

	// 创建更新器，配置资源过滤器
	// go-selfupdate 需要匹配格式: {name}_{os}_{arch}.{ext}
	updater, err := selfupdate.NewUpdater(selfupdate.Config{
		Source: u.source,
		// 匹配资源名称: wails-demo_windows_amd64.zip 或 wails-demo-windows-amd64.zip
		Filters: []string{
			fmt.Sprintf(`wails-demo[_-]%s[_-]%s\.zip`, runtime.GOOS, runtime.GOARCH),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("创建更新器失败: %w", err)
	}

	// 检查更新
	repoPath := fmt.Sprintf("%s/%s", u.owner, u.repo)
	latest, found, err := updater.DetectLatest(ctx, selfupdate.ParseSlug(repoPath))
	if err != nil {
		info.DebugInfo = fmt.Sprintf("检查更新出错: %v", err)
		return nil, fmt.Errorf("检查更新失败: %w", err)
	}

	info.DebugInfo = fmt.Sprintf("found=%v, repo=%s, currentVersion=%s", found, repoPath, currentVersion)

	if !found {
		info.LatestVersion = Version
		info.DebugInfo += ", 未找到匹配的 Release 资源"
		return info, nil
	}

	u.latest = latest
	info.LatestVersion = "v" + latest.Version()
	info.ReleaseNotes = latest.ReleaseNotes
	info.ReleaseURL = latest.URL
	info.DebugInfo += fmt.Sprintf(", latestVersion=%s", latest.Version())

	// 构建安装包下载 URL
	if u.mode == ModeInstaller {
		info.InstallerURL = fmt.Sprintf(
			"https://github.com/%s/%s/releases/download/v%s/wails-demo-setup-%s.exe",
			u.owner, u.repo, latest.Version(), latest.Version(),
		)
	}

	// 比较版本
	if latest.GreaterThan(currentVersion) {
		info.Available = true
	}

	return info, nil
}

// DownloadAndApplyUpdate 下载并应用更新
func (u *Updater) DownloadAndApplyUpdate(ctx context.Context) (*UpdateProgress, error) {
	if u.latest == nil {
		return nil, fmt.Errorf("请先检查更新")
	}

	if u.mode == ModeInstaller {
		return u.downloadAndRunInstaller(ctx)
	}

	return u.downloadAndReplaceBinary(ctx)
}

// downloadAndReplaceBinary 便携版：下载并替换二进制文件
func (u *Updater) downloadAndReplaceBinary(ctx context.Context) (*UpdateProgress, error) {
	// 创建更新器
	updater, err := selfupdate.NewUpdater(selfupdate.Config{
		Source: u.source,
		Filters: []string{
			fmt.Sprintf(`wails-demo[_-]%s[_-]%s\.zip`, runtime.GOOS, runtime.GOARCH),
		},
	})
	if err != nil {
		return &UpdateProgress{
			Status:  "error",
			Message: fmt.Sprintf("创建更新器失败: %v", err),
		}, err
	}

	// 获取当前可执行文件路径
	exe, err := selfupdate.ExecutablePath()
	if err != nil {
		return &UpdateProgress{
			Status:  "error",
			Message: fmt.Sprintf("获取可执行文件路径失败: %v", err),
		}, err
	}

	// 应用更新
	if err := updater.UpdateTo(ctx, u.latest, exe); err != nil {
		return &UpdateProgress{
			Status:  "error",
			Message: fmt.Sprintf("应用更新失败: %v", err),
		}, err
	}

	return &UpdateProgress{
		Status:      "ready",
		Progress:    100,
		Message:     "更新已就绪，请重启应用以完成安装",
		NeedRestart: true,
	}, nil
}

// downloadAndRunInstaller 安装版：下载新安装包并运行
func (u *Updater) downloadAndRunInstaller(ctx context.Context) (*UpdateProgress, error) {
	progress := &UpdateProgress{
		Status:   "downloading",
		Progress: 0,
		Message:  "正在下载安装包...",
	}

	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "wails-demo-update")
	if err != nil {
		progress.Status = "error"
		progress.Message = fmt.Sprintf("创建临时目录失败: %v", err)
		return progress, err
	}
	u.tempDir = tempDir

	// 构建安装包 URL
	installerURL := fmt.Sprintf(
		"https://github.com/%s/%s/releases/download/v%s/wails-demo-setup-%s.exe",
		u.owner, u.repo, u.latest.Version(), u.latest.Version(),
	)

	// 下载安装包
	installerPath := filepath.Join(tempDir, fmt.Sprintf("wails-demo-setup-%s.exe", u.latest.Version()))
	if err := u.downloadFile(ctx, installerURL, installerPath); err != nil {
		progress.Status = "error"
		progress.Message = fmt.Sprintf("下载安装包失败: %v", err)
		return progress, err
	}

	progress.Progress = 90
	progress.Message = "正在启动安装程序..."

	// 运行安装程序（静默模式）
	cmd := exec.Command(installerPath, "/S")
	if err := cmd.Start(); err != nil {
		progress.Status = "error"
		progress.Message = fmt.Sprintf("启动安装程序失败: %v", err)
		return progress, err
	}

	progress.Status = "ready"
	progress.Progress = 100
	progress.Message = "安装程序已启动，当前应用将自动关闭"
	progress.NeedRestart = true

	return progress, nil
}

// downloadFile 下载文件
func (u *Updater) downloadFile(ctx context.Context, url, filepath string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载失败，状态码: %d", resp.StatusCode)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// GetPlatformInfo 获取平台信息
func (u *Updater) GetPlatformInfo() string {
	return fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
}

// Cleanup 清理临时文件
func (u *Updater) Cleanup() {
	if u.tempDir != "" {
		os.RemoveAll(u.tempDir)
	}
}

// 确保 regexp 被使用（用于 Filters 模式）
var _ = regexp.Compile
