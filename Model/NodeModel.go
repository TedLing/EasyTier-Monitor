package Model

type Node struct {
	VirtualIP     string   `json:"virtual_ip"`
	Hostname      string   `json:"hostname"`
	ProxyCIDRs    string   `json:"proxy_cidrs"`
	PeerID        string   `json:"peer_id"`
	PublicIPv4    string   `json:"public_ipv4"`
	PublicIPv6    string   `json:"public_ipv6"`
	UDPStunType   string   `json:"udp_stun_type"`
	InterfaceIPv4 string   `json:"interface_ipv4"`
	InterfaceIPv6 string   `json:"interface_ipv6"`
	Listeners     []string `json:"listeners"` // 使用切片存储不定数量的 Listener
}
