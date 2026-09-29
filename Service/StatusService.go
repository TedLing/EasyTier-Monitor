package Service

import (
	"EasyTier-Monitor/Model"
	"EasyTier-Monitor/Tools"
	"encoding/json"
	"strconv"
	"strings"
)

// GetStatus 聚合 CLI node/peer/connector 数据，供状态页使用
func GetStatus() (Model.Status, error) {
	var status Model.Status

	// node 命令同时用于判断 core 是否运行
	res, err := Tools.RunCmd(Tools.AppConfig.CLI.Path, Tools.CLIArgs("-o", "json", "node")...)
	if err != nil {
		// CLI 无法连接 core，视为未运行
		Tools.AppLogger.Debug("获取node失败，core可能未运行: %v", err)
		return status, nil
	}

	var node Model.NodeNew
	if err := json.Unmarshal([]byte(res), &node); err != nil {
		Tools.AppLogger.Error("解析status node数据失败: %v", err)
		return status, nil
	}

	status.Running = true
	status.Version = node.Version
	status.Hostname = node.Hostname
	status.VirtualIP = node.IPv4Addr
	status.ProxyCIDRs = node.ProxyCidrs
	status.Listeners = node.Listeners

	// peer 信息（尽力获取，失败不影响整体）
	if peerRes, perr := Tools.RunCmd(Tools.AppConfig.CLI.Path, Tools.CLIArgs("-o", "json", "peer")...); perr == nil {
		var peers []Model.Peer
		if json.Unmarshal([]byte(peerRes), &peers) == nil {
			status.PeerCount = len(peers)
			for _, p := range peers {
				status.TotalRx += parseBytes(p.RxBytes)
				status.TotalTx += parseBytes(p.TxBytes)
			}
		}
	}

	// connector 数量（尽力获取）
	if connRes, cerr := Tools.RunCmd(Tools.AppConfig.CLI.Path, Tools.CLIArgs("-o", "json", "connector")...); cerr == nil {
		var connectors []Model.Connector
		if json.Unmarshal([]byte(connRes), &connectors) == nil {
			for _, c := range connectors {
				if c.Status == 0 {
					status.Connector++
				}
			}
		}
	}

	return status, nil
}

// parseBytes 将 CLI 返回的字节字符串转为 int64。
// JSON 输出中的 rx_bytes/tx_bytes 是人类可读格式（如 "20.57 kB"、"131.50 kB"、"-"），需按单位换算。
func parseBytes(s string) int64 {
	s = strings.TrimSpace(s)
	if s == "" || s == "-" {
		return 0
	}

	// 纯数字直接解析
	if v, err := strconv.ParseInt(s, 10, 64); err == nil {
		return v
	}

	// 拆分数值与单位
	parts := strings.Fields(s)
	if len(parts) == 0 {
		return 0
	}
	v, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0
	}

	unit := ""
	if len(parts) > 1 {
		unit = strings.ToUpper(parts[1])
	}
	multiplier := float64(1)
	switch {
	case strings.HasPrefix(unit, "K"):
		multiplier = 1 << 10
	case strings.HasPrefix(unit, "M"):
		multiplier = 1 << 20
	case strings.HasPrefix(unit, "G"):
		multiplier = 1 << 30
	case strings.HasPrefix(unit, "T"):
		multiplier = 1 << 40
	}
	return int64(v * multiplier)
}
