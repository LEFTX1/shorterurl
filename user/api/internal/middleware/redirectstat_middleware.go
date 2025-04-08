package middleware

import (
	"context"
	"net/http"

	"shorterurl/user/api/internal/types"
	"shorterurl/user/api/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

// 定义上下文键
const (
	RedirectStatsKey = "redirect_stats"
)

// redirectStatMiddleware 处理短链接跳转和统计数据收集
type redirectStatMiddleware struct{}

// NewRedirectStatMiddleware 创建一个新的RedirectStatMiddleware实例
func NewRedirectStatMiddleware() rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return (&redirectStatMiddleware{}).Handle(next)
	}
}

// Handle 是中间件的核心处理函数，对每个HTTP请求进行统计信息收集
func (m *redirectStatMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取IP地址
		ip := util.GetClientIP(r)

		// 获取地理位置信息
		locale, err := util.GetIPLocation(ip)
		if err != nil {
			logx.Errorf("获取IP地理位置失败: %v", err)
			locale = "未知"
		}

		// 获取短链接URI
		shortUri := r.URL.Path
		// 去除前导斜杠
		if len(shortUri) > 0 && shortUri[0] == '/' {
			shortUri = shortUri[1:]
		}

		// 创建统计信息
		stats := &types.RedirectStats{
			Ip:        ip,
			ShortUri:  shortUri,
			UserAgent: r.UserAgent(),
			Browser:   "未知",
			Os:        "未知",
			Device:    "未知",
			Network:   util.DetectNetworkType(ip, r.UserAgent()),
			Locale:    locale,
		}

		// 解析User-Agent
		if ua := r.UserAgent(); ua != "" {
			uaInfo := util.ParseUserAgent(ua)
			stats.Browser = uaInfo.Browser
			stats.Os = uaInfo.OS
			stats.Device = uaInfo.Device
			stats.DeviceModel = uaInfo.DeviceModel
		}

		// 将统计信息存储到context中
		ctx := context.WithValue(r.Context(), RedirectStatsKey, stats)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
