package util

import (
	"net/url"
	"strings"
)

// ExtractDomain 从URL中提取域名
func ExtractDomain(originUrl string) string {
	if originUrl == "" {
		return ""
	}

	// 确保URL有协议前缀
	if !strings.HasPrefix(originUrl, "http://") && !strings.HasPrefix(originUrl, "https://") {
		originUrl = "http://" + originUrl
	}

	// 解析URL
	parsedURL, err := url.Parse(originUrl)
	if err != nil {
		return ""
	}

	// 获取主机部分
	host := parsedURL.Host

	// 移除端口号部分
	if strings.Contains(host, ":") {
		host = strings.Split(host, ":")[0]
	}

	// 如果是子域名，返回主域名
	parts := strings.Split(host, ".")
	if len(parts) > 2 {
		return strings.Join(parts[len(parts)-2:], ".")
	}

	return host
}
