package Model

// Status 状态页聚合数据
type Status struct {
	Running     bool     `json:"running"`      // easytier-core 是否运行
	Version     string   `json:"version"`      // 内核版本
	Hostname    string   `json:"hostname"`     // 主机名
	VirtualIP   string   `json:"virtual_ip"`   // 虚拟网卡 IPv4
	ProxyCIDRs  []string `json:"proxy_cidrs"`  // 代理 CIDR
	Listeners   []string `json:"listeners"`    // 监听地址
	PeerCount   int      `json:"peer_count"`   // 在线节点数
	TotalRx     int64    `json:"total_rx"`     // 累计接收字节（所有节点汇总）
	TotalTx     int64    `json:"total_tx"`     // 累计发送字节（所有节点汇总）
	Connector   int      `json:"connector"`    // 已连接的服务器节点数
}
