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

	return host
}
