package Model

// ConnectionInfo 连接信息页聚合数据，各字段为 easytier-cli 子命令的原始文本输出
type ConnectionInfo struct {
	Node           string `json:"node"`
	Peer           string `json:"peer"`
	Route          string `json:"route"`
	Connector      string `json:"connector"`
	Stun           string `json:"stun"`
	PeerCenter     string `json:"peer_center"`
	VPNPortal      string `json:"vpn_portal"`
	Proxy          string `json:"proxy"`
	ACL            string `json:"acl"`
	MappedListener string `json:"mapped_listener"`
	Stats          string `json:"stats"`
	Cmdline        string `json:"cmdline"`
	ConfigFile     string `json:"config_file"`
}
