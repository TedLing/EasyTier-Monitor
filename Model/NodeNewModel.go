package Model

type NodeNew struct {
	PeerID      uint32      `json:"peer_id"`
	IPv4Addr    string      `json:"ipv4_addr"`
	ProxyCidrs  []string    `json:"proxy_cidrs"`
	Hostname    string      `json:"hostname"`
	StunInfo    StunInfo    `json:"stun_info"`
	InstID      string      `json:"inst_id"`
	Listeners   []string    `json:"listeners"`
	Config      string      `json:"config"`
	Version     string      `json:"version"`
	FeatureFlag FeatureFlag `json:"feature_flag"`
	IPList      IPList      `json:"ip_list"`
}

type StunInfo struct {
	UDPNatType     int      `json:"udp_nat_type"`
	TCPNatType     int      `json:"tcp_nat_type"`
	LastUpdateTime int64    `json:"last_update_time"`
	PublicIP       []string `json:"public_ip"`
	MinPort        int      `json:"min_port"`
	MaxPort        int      `json:"max_port"`
}

type FeatureFlag struct {
	IsPublicServer bool `json:"is_public_server"`
	AvoidRelayData bool `json:"avoid_relay_data"`
	KCPInput       bool `json:"kcp_input"`
	NoRelayKCP     bool `json:"no_relay_kcp"`
}

type IPList struct {
	PublicIPv4     IPAddr     `json:"public_ipv4"`
	InterfaceIPv4s []IPAddr   `json:"interface_ipv4s"`
	PublicIPv6     IPv6Addr   `json:"public_ipv6"`
	InterfaceIPv6s []IPv6Addr `json:"interface_ipv6s"`
	Listeners      []string   `json:"listeners"`
}

type IPAddr struct {
	Addr uint32 `json:"addr"`
}

type IPv6Addr struct {
	Part1 uint32 `json:"part1"`
	Part2 uint32 `json:"part2"`
	Part3 uint32 `json:"part3"`
	Part4 uint32 `json:"part4"`
}
