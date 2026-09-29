// 定义API响应的基本结构
export interface ApiResponse<T = any> {
  code: number;
  msg: string;
  data: T;
}

// 节点信息接口
export interface NodeInfo {
  virtual_ip: string;
  hostname: string;
  proxy_cidrs: string;
  public_ipv4: string;
  public_ipv6: string;
  udp_stun_type?: string;
  interface_ipv4: string;
  interface_ipv6: string;
  listeners: string;
}

// 服务器连接信息接口
export interface Connector {
  url: string;
  status: string;
}

// 对等节点信息接口
export interface Peer {
  id: string;
  ipv4: string;
  hostname: string;
  cost: string;
  lat_ms: number;
  loss_rate: number;
  rx_bytes: number;
  tx_bytes: number;
  tunnel_proto: string;
  nat_type: string;
  version: string;
}
