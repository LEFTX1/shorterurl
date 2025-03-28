package util

import (
	"fmt"
	"net/url"
	"strings"
)

// GetFavicon 获取网站favicon图标URL
// 根据原始URL提取域名，并生成favicon的URL
func GetFavicon(originUrl string) string {
	if originUrl == "" {
		return ""
	}

	// 解析URL获取域名
	parsedUrl, err := url.Parse(originUrl)
	if err != nil || parsedUrl.Host == "" {
		// 如果URL解析失败，尝试修复URL
		if !strings.HasPrefix(originUrl, "http://") && !strings.HasPrefix(originUrl, "https://") {
			originUrl = "http://" + originUrl
			parsedUrl, err = url.Parse(originUrl)
			if err != nil || parsedUrl.Host == "" {
				return ""
			}
		} else {
			return ""
		}
	}

	// 获取域名
	domain := parsedUrl.Host

	// 使用Google的favicon服务获取图标
	// 也可以替换为其他favicon服务
	return fmt.Sprintf("https://www.google.com/s2/favicons?domain=%s", domain)
}
