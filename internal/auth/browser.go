package auth

import (
	"os"
	"os/exec"
	"runtime"
)

// openBrowser 跨平台打开浏览器
func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	go cmd.Wait() // 回收子进程，避免僵尸进程
	return nil
}

// isLocalEnvironment 检测是否为本地桌面环境
func isLocalEnvironment() bool {
	if runtime.GOOS == "darwin" || runtime.GOOS == "windows" {
		return true
	}
	// Linux：检查 DISPLAY 或 WAYLAND_DISPLAY
	if os.Getenv("DISPLAY") != "" || os.Getenv("WAYLAND_DISPLAY") != "" {
		return true
	}
	return false
}
