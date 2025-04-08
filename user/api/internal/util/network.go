package util

import (
	"net"
	"strings"
)

// 移动网络IP段
var mobileIPRanges = []string{
	"10.0.0.0/8",     // 私有网络
	"172.16.0.0/12",  // 私有网络
	"192.168.0.0/16", // 私有网络
	"100.64.0.0/10",  // 共享地址空间
	"169.254.0.0/16", // 链路本地
}

// 移动设备关键词
var mobileKeywords = []string{
	"mobile",
	"android",
	"iphone",
	"ipad",
	"ipod",
	"windows phone",
}

// DetectNetworkType 推测网络类型
func DetectNetworkType(ip string, userAgent string) string {
	// 1. 检查是否是本地IP
	if isLocalIP(ip) {
		return "local"
	}

	// 2. 检查是否是移动设备
	if isMobileDevice(userAgent) {
		// 移动设备更可能使用移动网络
		return "mobile"
	}

	// 3. 默认返回wifi
	return "wifi"
}

// isLocalIP 检查是否是本地IP
func isLocalIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	for _, cidr := range mobileIPRanges {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		if network.Contains(parsedIP) {
			return true
		}
	}

	return false
}

// isMobileDevice 检查是否是移动设备
func isMobileDevice(userAgent string) bool {
	ua := strings.ToLower(userAgent)
	for _, keyword := range mobileKeywords {
		if strings.Contains(ua, keyword) {
			return true
		}
	}
	return false
}
