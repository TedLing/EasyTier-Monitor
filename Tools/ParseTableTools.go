package Tools

import (
	"EasyWeb/Model"
	"encoding/json"
	"fmt"
	"strings"
)

// ParsePeerTableToJSON 节点解析
func ParsePeerTableToJSON(tableStr string) (string, error) {
	// 按行分割
	lines := strings.Split(strings.TrimSpace(tableStr), "\n")
	if len(lines) < 3 {
		return "", fmt.Errorf("invalid table format: too few lines")
	}

	// 提取表头（第二行，忽略边框）
	headerLine := strings.TrimSpace(lines[1])
	// 分割表头并清理首尾空字段
	headers := strings.Split(headerLine, "│")
	var cleanHeaders []string
	for _, h := range headers {
		h = strings.TrimSpace(h)
		if h != "" {
			cleanHeaders = append(cleanHeaders, h)
		}
	}
	if len(cleanHeaders) < 11 {
		return "", fmt.Errorf("invalid table format: too few columns (%d)", len(cleanHeaders))
	}

	// 解析数据行（从第三行到倒数第二行，忽略分隔线和边框）
	var rows []map[string]string
	for i := 3; i < len(lines)-1; i++ {
		if strings.Contains(lines[i], "├") || strings.Contains(lines[i], "└") {
			continue // 跳过分隔线
		}
		fields := strings.Split(strings.TrimSpace(lines[i]), "│")
		// 清理首尾空字段
		var cleanFields []string
		for _, f := range fields {
			f = strings.TrimSpace(f)
			if f != "" {
				cleanFields = append(cleanFields, f)
			}
		}
		if len(cleanFields) != len(cleanHeaders) {
			continue // 跳过格式错误的行
		}

		// 创建 map 表示一行
		row := make(map[string]string)
		for j, field := range cleanFields {
			row[cleanHeaders[j]] = field
		}
		rows = append(rows, row)
	}

	// 转换为 JSON
	jsonData, err := json.MarshalIndent(rows, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %v", err)
	}
	return string(jsonData), nil
}

// ParseNodeToModel  node节点解析
func ParseNodeToModel(output string) (Model.Node, error) {
	var node Model.Node
	node.Listeners = []string{} // 初始化 Listeners 切片
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		// 跳过分隔线和空行
		if !strings.Contains(line, "│") || strings.Contains(line, "────") {
			continue
		}

		// 分割字段
		fields := strings.Split(line, "│")
		if len(fields) < 3 {
			continue
		}

		// 清理键和值
		key := strings.TrimSpace(fields[1])
		value := strings.TrimSpace(fields[2])

		// 根据键分配值
		switch {
		case key == "Virtual IP":
			node.VirtualIP = value
		case key == "Hostname":
			node.Hostname = value
		case key == "Proxy CIDRs":
			node.ProxyCIDRs = value
		case key == "Peer ID":
			node.PeerID = value
		case key == "Public IPv4":
			node.PublicIPv4 = value
		case key == "UDP Stun Type":
			node.UDPStunType = value
		case key == "Interface IPv4":
			node.InterfaceIPv4 = value
		case key == "Interface IPv6":
			node.InterfaceIPv6 = value
		case strings.HasPrefix(key, "Listener"): // 动态处理所有 Listener
			node.Listeners = append(node.Listeners, value)
		}
	}

	return node, nil
}
