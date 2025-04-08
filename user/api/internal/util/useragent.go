package util

import (
	"strings"
)

// UserAgentInfo 存储解析后的User-Agent信息
type UserAgentInfo struct {
	Browser     string // 浏览器名称
	OS          string // 操作系统
	Device      string // 设备类型
	DeviceModel string // 设备型号
}

// ParseUserAgent 解析User-Agent字符串
func ParseUserAgent(userAgent string) *UserAgentInfo {
	info := &UserAgentInfo{
		Browser:     "unknown",
		OS:          "unknown",
		Device:      "unknown",
		DeviceModel: "unknown",
	}

	if userAgent == "" {
		return info
	}

	ua := strings.ToLower(userAgent)

	// 解析浏览器信息
	switch {
	case strings.Contains(ua, "chrome"):
		info.Browser = "Chrome"
	case strings.Contains(ua, "firefox"):
		info.Browser = "Firefox"
	case strings.Contains(ua, "safari"):
		info.Browser = "Safari"
	case strings.Contains(ua, "edge"):
		info.Browser = "Edge"
	case strings.Contains(ua, "opera"):
		info.Browser = "Opera"
	case strings.Contains(ua, "msie") || strings.Contains(ua, "trident"):
		info.Browser = "Internet Explorer"
	}

	// 解析操作系统信息
	switch {
	case strings.Contains(ua, "windows"):
		info.OS = "Windows"
	case strings.Contains(ua, "macintosh") || strings.Contains(ua, "mac os x"):
		info.OS = "macOS"
	case strings.Contains(ua, "linux"):
		info.OS = "Linux"
	case strings.Contains(ua, "android"):
		info.OS = "Android"
		info.Device = "Mobile"
		info.DeviceModel = extractAndroidModel(ua)
	case strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad") || strings.Contains(ua, "ipod"):
		info.OS = "iOS"
		if strings.Contains(ua, "iphone") {
			info.Device = "iPhone"
		} else if strings.Contains(ua, "ipad") {
			info.Device = "iPad"
		} else {
			info.Device = "iPod"
		}
	}

	// 如果没有检测到移动设备，则判断为桌面设备
	if info.Device == "unknown" {
		info.Device = "Desktop"
		if info.OS == "Windows" {
			info.DeviceModel = "Windows"
		} else if info.OS == "macOS" {
			info.DeviceModel = "Mac"
		} else if info.OS == "Linux" {
			info.DeviceModel = "Linux"
		}
	}

	return info
}

// extractAndroidModel 提取Android设备型号
func extractAndroidModel(ua string) string {
	// 尝试从User-Agent中提取设备型号
	start := strings.Index(ua, "build/")
	if start == -1 {
		return "unknown"
	}
	start = strings.Index(ua[start:], ";") + start + 1
	if start == 0 {
		return "unknown"
	}
	end := strings.Index(ua[start:], " ")
	if end == -1 {
		return ua[start:]
	}
	return ua[start : start+end]
}
