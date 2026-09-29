//go:build linux

package Service

import (
	"os"
	"strings"
)

// getCoreCmdline 读取 easytier-core 进程启动参数（对应官方 /proc/<pid>/cmdline）
func getCoreCmdline() string {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return ""
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		data, err := os.ReadFile("/proc/" + entry.Name() + "/cmdline")
		if err != nil {
			continue
		}
		if strings.Contains(string(data), "easytier-core") {
			return strings.ReplaceAll(string(data), "\x00", " ")
		}
	}
	return ""
}

// readCoreConfig 读取 core 配置文件
func readCoreConfig() string {
	data, err := os.ReadFile("/etc/easytier/config.toml")
	if err != nil {
		return ""
	}
	return string(data)
}
