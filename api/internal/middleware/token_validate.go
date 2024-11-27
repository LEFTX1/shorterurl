package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest/httpx"

	"shorterurl/api/internal/config"
	"shorterurl/api/internal/types"
)

// TokenValidateMiddleware 是一个中间件结构体，用于验证请求中的 token
type TokenValidateMiddleware struct {
	Config     *config.TokenValidateConfig // 中间件配置，包含白名单等信息
	RedisCache *redis.Redis                // Redis 客户端，用于存储和验证 token
}

// NewTokenValidateMiddleware 创建一个新的 TokenValidateMiddleware 实例
func NewTokenValidateMiddleware(config *config.TokenValidateConfig, redisCache *redis.Redis) *TokenValidateMiddleware {
	return &TokenValidateMiddleware{
		Config:     config,
		RedisCache: redisCache,
	}
}

// Handle 是中间件的核心处理函数，对每个 HTTP 请求进行 token 验证
func (m *TokenValidateMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 记录请求信息
		logx.Infof("网关收到请求 - 路径: %s, 方法: %s", r.URL.Path, r.Method)

		// 如果请求路径在白名单中，直接放行
		if m.isPathInWhiteList(r.URL.Path, r.Method) {
			logx.Infof("请求在白名单中，直接放行 - 路径: %s", r.URL.Path)
			next(w, r)
			return
		}

		// 从请求头中获取用户凭证
		username := r.Header.Get("username")
		token := r.Header.Get("token")

		// 验证用户凭证是否存在
		if username == "" || token == "" {
			// 记录凭证缺失错误
			logx.WithContext(r.Context()).Errorf("用户名或令牌为空 - 路径: %s, 用户名: %s", r.URL.Path, username)
			// 返回未授权错误响应
			httpx.WriteJson(w, http.StatusUnauthorized, &types.GatewayErrorResult{
				Status:  http.StatusUnauthorized,
				Message: "用户名或令牌不能为空",
			})
			return
		}

		// 从 Redis 中获取用户信息，验证 token
		userInfoStr, err := m.RedisCache.Hget("short-link:login:"+username, token)
		if err != nil || userInfoStr == "" {
			// 记录 token 无效错误
			logx.WithContext(r.Context()).Errorf("令牌在 Redis 中未找到 - 路径: %s, 用户名: %s", r.URL.Path, username)
			// 返回未授权错误响应
			httpx.WriteJson(w, http.StatusUnauthorized, &types.GatewayErrorResult{
				Status:  http.StatusUnauthorized,
				Message: "无效的令牌",
			})
			return
		}

		// 解析存储在 Redis 中的用户信息
		var userInfo map[string]interface{}
		if err := json.Unmarshal([]byte(userInfoStr), &userInfo); err != nil {
			// 记录用户信息解析错误
			logx.WithContext(r.Context()).Errorf("解析用户信息失败 - 路径: %s, 用户名: %s", r.URL.Path, username)
			// 返回未授权错误响应
			httpx.WriteJson(w, http.StatusUnauthorized, &types.GatewayErrorResult{
				Status:  http.StatusUnauthorized,
				Message: "解析用户信息失败",
			})
			return
		}

		// 创建用户上下文信息对象
		ctxUserInfo := &types.UserInfo{
			ID:       userInfo["id"].(string),       // 用户ID
			Username: username,                      // 用户名
			RealName: userInfo["realName"].(string), // 真实姓名
		}

		// 将用户信息添加到请求上下文中
		ctx := context.WithValue(r.Context(), types.UserContextKey, ctxUserInfo)

		// 使用更新后的上下文创建新的请求对象
		r = r.WithContext(ctx)

		// 调用下一个处理函数
		next(w, r)
	}
}

// isPathInWhiteList 检查请求路径是否在白名单中
func (m *TokenValidateMiddleware) isPathInWhiteList(path, method string) bool {
	// 检查是否是用户注册接口（特殊白名单）
	if path == "/api/short-link/admin/v1/user" && method == "POST" {
		return true
	}

	// 检查是否在配置的白名单路径中
	for _, whitePath := range m.Config.WhitePathList {
		if strings.HasPrefix(path, whitePath) {
			return true
		}
	}

	return false
}
