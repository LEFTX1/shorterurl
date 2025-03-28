package util

import (
	"time"
)

// 常量定义
const (
	// 24小时的毫秒数
	DayMilliseconds = 24 * 60 * 60 * 1000
	// 永久有效类型
	ValidDateTypePermanent = 0
	// 自定义有效期类型
	ValidDateTypeCustom = 1
)

// GetLinkCacheValidTime 计算链接缓存的有效时间（毫秒）
// 如果是永久有效，则返回30天的毫秒数
// 如果是自定义有效期，则返回从现在到有效期的毫秒数
func GetLinkCacheValidTime(validDate time.Time) int64 {
	// 当前时间
	now := time.Now()

	// 如果有效期为空或者过去的时间，则返回30天缓存时间
	if validDate.IsZero() || validDate.Before(now) {
		return 30 * DayMilliseconds // 30天
	}

	// 计算从now到validDate的毫秒数
	return validDate.UnixMilli() - now.UnixMilli()
}

// IsValidLink 判断链接是否有效
// 根据短链接的有效期类型和有效期判断链接是否过期
func IsValidLink(validDateType int, validDate time.Time) bool {
	// 如果是永久有效类型，直接返回true
	if validDateType == ValidDateTypePermanent {
		return true
	}

	// 如果是自定义有效期，判断当前时间是否超过有效期
	if validDateType == ValidDateTypeCustom {
		return time.Now().Before(validDate)
	}

	// 默认返回false
	return false
}
