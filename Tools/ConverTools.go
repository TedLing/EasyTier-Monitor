package Tools

import (
	"EasyTier-Monitor/Model"
	"encoding/binary"
	"net"
	"strings"
)

func ToIPv4(ip Model.IPAddr) string {
	// 将 uint32 转换为 4 字节的 IPv4 地址
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, ip.Addr) // 假设 addr 是大端序
	return net.IP(b).String()
}

func ToIPv4List(ips []Model.IPAddr) string {

	var ipsStr []string
	for _, ip := range ips {
		ipv4 := ToIPv4(ip)
		ipsStr = append(ipsStr, ipv4)
	}

	return strings.Join(ipsStr, ",")

}

func ToIPv6(ip Model.IPv6Addr) string {
	// 创建 16 字节切片（IPv6 地址为 128 位）
	b := make([]byte, 16)
	// 将每个 uint32 写入字节切片（大端序）
	binary.BigEndian.PutUint32(b[0:4], ip.Part1)
	binary.BigEndian.PutUint32(b[4:8], ip.Part2)
	binary.BigEndian.PutUint32(b[8:12], ip.Part3)
	binary.BigEndian.PutUint32(b[12:16], ip.Part4)
	// 转换为 IPv6 地址字符串
	return net.IP(b).String()
}

// ToIPv6List 批量转换
func ToIPv6List(ips []Model.IPv6Addr) string {

	var ipsStr []string
	for _, ip := range ips {
		ipv6 := ToIPv6(ip)
		ipsStr = append(ipsStr, ipv6)
	}

	return strings.Join(ipsStr, ",")

}
