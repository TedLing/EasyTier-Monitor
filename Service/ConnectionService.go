package Service

import (
	"EasyTier-Monitor/Model"
	"EasyTier-Monitor/Tools"
	"strings"
)

// GetConnection 聚合 easytier-cli 各子命令的原始文本，供连接信息页解析展示。
// 复刻官方 luci 插件 act_conninfo 的数据来源。
func GetConnection() (Model.ConnectionInfo, error) {
	var info Model.ConnectionInfo

	// node 用于判断 core 是否运行：连不上时 CLI 输出含连接拒绝/无实例等错误
	nodeOut, err := Tools.RunCmdCombined(Tools.AppConfig.CLI.Path, Tools.CLIArgs("node")...)
	if err != nil {
		return info, err
	}

	// core 未运行时所有字段填充统一错误文案
	if !isCoreRunning(nodeOut) {
		errMsg := "Error: Program not running! Please start the program and refresh"
		info.Node = errMsg
		info.Peer = errMsg
		info.Route = errMsg
		info.Connector = errMsg
		info.Stun = errMsg
		info.PeerCenter = errMsg
		info.VPNPortal = errMsg
		info.Proxy = errMsg
		info.ACL = errMsg
		info.MappedListener = errMsg
		info.Stats = errMsg
		info.Cmdline = errMsg
		return info, nil
	}

	info.Node = nodeOut
	info.Peer = cliRaw("peer")
	info.Route = cliRaw("route")
	info.Connector = cliRaw("connector")
	info.Stun = cliRaw("stun")
	info.PeerCenter = cliRaw("peer-center")
	info.VPNPortal = cliRaw("vpn-portal")
	info.Proxy = cliRaw("proxy")
	info.ACL = cliRaw("acl", "stats")
	info.MappedListener = cliRaw("mapped-listener")
	info.Stats = cliRaw("stats")

	info.Cmdline = getCoreCmdline()
	if strings.Contains(info.Cmdline, "--config-file") {
		info.ConfigFile = readCoreConfig()
	}

	return info, nil
}

// cliRaw 执行单个子命令并返回原始文本
func cliRaw(args ...string) string {
	out, err := Tools.RunCmdCombined(Tools.AppConfig.CLI.Path, Tools.CLIArgs(args...)...)
	if err != nil {
		Tools.AppLogger.Error("执行连接信息命令 %v 失败: %v", args, err)
		return ""
	}
	return out
}

// isCoreRunning 根据 node 命令输出判断 core 是否在运行
func isCoreRunning(nodeOut string) bool {
	lower := strings.ToLower(nodeOut)
	for _, marker := range []string{
		"connection refused",
		"no running instances found",
		"os error",
		"failed to connect",
	} {
		if strings.Contains(lower, marker) {
			return false
		}
	}
	return strings.Contains(nodeOut, "|")
}
