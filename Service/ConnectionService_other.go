//go:build !linux

package Service

// getCoreCmdline Windows/macOS 下无 /proc，无法获取 core 启动参数
func getCoreCmdline() string {
	return ""
}

// readCoreConfig 非官方支持平台不读取固定路径配置
func readCoreConfig() string {
	return ""
}
