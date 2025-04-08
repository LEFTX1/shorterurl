package util

import (
	"net/http"
	"strings"
)

// GetClientIP 获取客户端真实IP地址
func GetClientIP(r *http.Request) string {
	// 尝试从X-Real-IP获取
	ip := r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	// 尝试从X-Forwarded-For获取
	ip = r.Header.Get("X-Forwarded-For")
	if ip != "" {
		// 取第一个IP
		ips := strings.Split(ip, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// 从RemoteAddr获取
	ip = r.RemoteAddr
	if ip != "" {
		// 去除端口号
		if idx := strings.LastIndex(ip, ":"); idx != -1 {
			ip = ip[:idx]
		}
		// 处理IPv6地址
		if strings.HasPrefix(ip, "[") && strings.HasSuffix(ip, "]") {
			ip = ip[1 : len(ip)-1]
		}
		return ip
	}

	return "unknown"
}
