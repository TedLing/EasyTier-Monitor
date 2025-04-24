package Model

type Peer struct {
	IPv4        string `json:"ipv4"`
	Hostname    string `json:"hostname"`
	Cost        string `json:"cost"`
	LatMs       string `json:"lat_ms"`
	LossRate    string `json:"loss_rate"`
	RxBytes     string `json:"rx_bytes"`
	TxBytes     string `json:"tx_bytes"`
	TunnelProto string `json:"tunnel_proto"`
	NatType     string `json:"nat_type"`
	ID          string `json:"id"`
	Version     string `json:"version"`
}
