package Tools

import (
	"net"
	"strings"
)

// ValidateIPv4 验证并格式化IPv4地址
func ValidateIPv4(ipStr string) string {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return ""
	}
	ipv4 := ip.To4()
	if ipv4 == nil {
		return ""
	}
	return ipv4.String()
}

// ValidateIPv6 验证并格式化IPv6地址
func ValidateIPv6(ipStr string) string {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return ""
	}
	ipv6 := ip.To16()
	if ipv6 == nil {
		return ""
	}
	return ipv6.String()
}

// ValidateIPv4List 验证并格式化IPv4地址列表
func ValidateIPv4List(ipList []string) string {
	var validIPs []string
	for _, ipStr := range ipList {
		if ipv4 := ValidateIPv4(ipStr); ipv4 != "" {
			validIPs = append(validIPs, ipv4)
		}
	}
	return strings.Join(validIPs, ",")
}

// ValidateIPv6List 验证并格式化IPv6地址列表
func ValidateIPv6List(ipList []string) string {
	var validIPs []string
	for _, ipStr := range ipList {
		if ipv6 := ValidateIPv6(ipStr); ipv6 != "" {
			validIPs = append(validIPs, ipv6)
		}
	}
	return strings.Join(validIPs, ",")
}
