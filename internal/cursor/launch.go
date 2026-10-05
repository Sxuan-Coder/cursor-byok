package cursor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	goruntime "runtime"
	"strings"

	"cursor/internal/logger"
)

// Launch 启动 Cursor 客户端，不等待进程退出。
func Launch() error {
	executablePath := resolveCursorExecutable()
	if executablePath == "" {
		return launchCursorFromSystem()
	}

	if err := exec.Command(executablePath).Start(); err != nil {
		return fmt.Errorf("启动 Cursor 失败 path=%s: %w", executablePath, err)
	}

	logger.Infof("launchCursor: started path=%s", executablePath)
	return nil
}

// resolveCursorExecutable 依次从环境变量、PATH 和常见安装目录里查找 Cursor 可执行文件。
func resolveCursorExecutable() string {
	for _, candidate := range cursorExecutableCandidates() {
		if isExecutableFile(candidate) {
			return candidate
		}
	}

	if path, err := exec.LookPath(cursorExecutableName()); err == nil {
		return path
	}

	return ""
}

func isExecutableFile(path string) bool {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return false
	}

	info, err := os.Stat(trimmed)
	if err != nil || info.IsDir() {
		return false
	}

	return true
}

func cursorExecutableName() string {
	if goruntime.GOOS == "windows" {
		return "Cursor.exe"
	}
	return "cursor"
}

// cursorExecutableCandidates 返回当前系统下 Cursor 的候选安装路径。
func cursorExecutableCandidates() []string {
	switch goruntime.GOOS {
	case "windows":
		return windowsCursorCandidates()
	case "darwin":
		return []string{
			"/Applications/Cursor.app/Contents/MacOS/Cursor",
			filepath.Join(userHomeDir(), "Applications/Cursor.app/Contents/MacOS/Cursor"),
		}
	default:
		return []string{
			"/usr/bin/cursor",
			"/usr/local/bin/cursor",
			"/opt/Cursor/cursor",
			"/var/lib/flatpak/exports/bin/com.cursor.Cursor",
			filepath.Join(userHomeDir(), ".local/bin/cursor"),
			filepath.Join(userHomeDir(), ".local/share/applications/cursor"),
		}
	}
}

func windowsCursorCandidates() []string {
	roots := []string{
		os.Getenv("LOCALAPPDATA"),
		os.Getenv("PROGRAMFILES"),
		os.Getenv("PROGRAMFILES(X86)"),
	}

	candidates := make([]string, 0, len(roots))
	for _, root := range roots {
		trimmed := strings.TrimSpace(root)
		if trimmed == "" {
			continue
		}

		for _, relativePath := range []string{
			filepath.Join("Programs", "cursor", "Cursor.exe"),
			filepath.Join("Programs", "Cursor", "Cursor.exe"),
			filepath.Join("cursor", "Cursor.exe"),
			filepath.Join("Cursor", "Cursor.exe"),
		} {
			candidates = append(candidates, filepath.Join(trimmed, relativePath))
		}
	}

	return candidates
}

// launchCursorFromSystem 在找不到安装路径时交给系统级启动方式兜底。
func launchCursorFromSystem() error {
	switch goruntime.GOOS {
	case "darwin":
		if err := exec.Command("open", "-a", "Cursor").Start(); err != nil {
			return fmt.Errorf("启动 Cursor 失败: %w", err)
		}
		logger.Infof("launchCursor: started by open -a Cursor")
		return nil
	case "windows":
		return fmt.Errorf("未找到 Cursor 安装位置，请确认 Cursor 已安装或手动启动")
	default:
		if err := exec.Command("cursor").Start(); err != nil {
			return fmt.Errorf("启动 Cursor 失败: %w", err)
		}
		logger.Infof("launchCursor: started by PATH lookup")
		return nil
	}
}

func userHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return home
}