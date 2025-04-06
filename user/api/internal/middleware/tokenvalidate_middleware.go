package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest/httpx"

	"shorterurl/user/api/internal/config"
	"shorterurl/user/api/internal/types"
	"shorterurl/user/api/internal/types/errorx"
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
		// 记录请求信息 - 增强日志
		logx.Infof("[TokenValidate] 收到请求 - 路径: %s, 方法: %s, 来源: %s", r.URL.Path, r.Method, r.RemoteAddr)

		// 打印所有请求头，用于调试
		logx.Infof("[TokenValidate] 请求头信息:")
		for name, values := range r.Header {
			logx.Infof("[TokenValidate] %s: %s", name, strings.Join(values, ", "))
		}

		// 如果请求路径在白名单中，直接放行
		if m.isPathInWhiteList(r.URL.Path, r.Method) {
			logx.Infof("[TokenValidate] 请求在白名单中，直接放行 - 路径: %s", r.URL.Path)
			next(w, r)
			return
		}

		// 从请求头中获取用户凭证
		username := r.Header.Get("username")
		token := r.Header.Get("token")

		// 详细记录凭证信息
		logx.Infof("[TokenValidate] 从请求头获取 - username: '%s', token: '%s'", username, token)

		// 验证用户凭证是否存在
		if username == "" || token == "" {
			// 记录凭证缺失错误
			logx.WithContext(r.Context()).Errorf("[TokenValidate] 凭证缺失 - 路径: %s, username: '%s', token: '%s'", r.URL.Path, username, token)
			// 返回未授权错误响应
			err := errorx.New(errorx.ClientError, "UNAUTHORIZED", "用户名或令牌不能为空")
			httpx.WriteJson(w, http.StatusUnauthorized, GatewayErrorResult{
				Status:  http.StatusUnauthorized,
				Message: err.Message,
			})
			return
		}

		// 从 Redis 中获取用户信息，验证 token
		// 使用配置的LoginKeyPrefix替代硬编码的"short-link:login:"
		loginKeyPrefix := "user:login:" // 默认值
		if m.Config != nil && m.Config.LoginKeyPrefix != "" {
			loginKeyPrefix = m.Config.LoginKeyPrefix
		}

		// 记录要查询的Redis键
		redisKey := loginKeyPrefix + username
		logx.Infof("[TokenValidate] 查询Redis - 键: '%s', token: '%s'", redisKey, token)

		// 添加Redis键搜索
		logx.Infof("[TokenValidate] 检查Redis中所有用户登录键:")
		userKeys, err := m.RedisCache.KeysCtx(r.Context(), "user:login:*")
		if err != nil {
			logx.Errorf("[TokenValidate] 获取Redis键失败: %v", err)
		} else {
			for _, key := range userKeys {
				logx.Infof("[TokenValidate] 找到用户登录键: %s", key)
				// 打印此键下的所有field和value
				values, err := m.RedisCache.HgetallCtx(r.Context(), key)
				if err != nil {
					logx.Errorf("[TokenValidate] 获取键 %s 的值失败: %v", key, err)
				} else {
					for field, value := range values {
						logx.Infof("[TokenValidate] 键: %s, field: %s, value长度: %d", key, field, len(value))
					}
				}
			}
		}

		userInfoStr, err := m.RedisCache.Hget(redisKey, token)

		// 记录Redis查询结果
		if err != nil {
			logx.Errorf("[TokenValidate] Redis查询错误 - 键: '%s', 错误: %v", redisKey, err)
		} else {
			logx.Infof("[TokenValidate] Redis查询结果 - 用户信息: '%s'", userInfoStr)
		}

		if err != nil || userInfoStr == "" {
			// 记录 token 无效错误
			logx.WithContext(r.Context()).Errorf("[TokenValidate] 无效令牌 - 路径: %s, 用户名: '%s', Redis错误: %v", r.URL.Path, username, err)
			// 返回未授权错误响应
			err := errorx.New(errorx.ClientError, "INVALID_TOKEN", "无效的令牌")
			httpx.WriteJson(w, http.StatusUnauthorized, GatewayErrorResult{
				Status:  http.StatusUnauthorized,
				Message: err.Message,
			})
			return
		}

		// 解析存储在 Redis 中的用户信息
		var userInfo map[string]interface{}
		if err := json.Unmarshal([]byte(userInfoStr), &userInfo); err != nil {
			// 记录用户信息解析错误
			logx.WithContext(r.Context()).Errorf("[TokenValidate] JSON解析失败 - 路径: %s, 用户名: '%s', 数据: '%s', 错误: %v",
				r.URL.Path, username, userInfoStr, err)
			// 返回未授权错误响应
			err := errorx.New(errorx.SystemError, "PARSE_ERROR", "解析用户信息失败")
			httpx.WriteJson(w, http.StatusUnauthorized, GatewayErrorResult{
				Status:  http.StatusUnauthorized,
				Message: err.Message,
			})
			return
		}

		// 记录解析后的用户信息
		logx.Infof("[TokenValidate] 用户信息解析成功: %+v", userInfo)

		// 检查必要的用户信息字段
		idVal, hasID := userInfo["id"]
		realNameVal, hasRealName := userInfo["realName"]

		if !hasID || !hasRealName {
			logx.Errorf("[TokenValidate] 用户信息缺少必要字段 - hasID: %v, hasRealName: %v", hasID, hasRealName)
			err := errorx.New(errorx.SystemError, "INVALID_USER_INFO", "用户信息不完整")
			httpx.WriteJson(w, http.StatusUnauthorized, GatewayErrorResult{
				Status:  http.StatusUnauthorized,
				Message: err.Message,
			})
			return
		}

		// 类型检查和转换
		idStr, idOK := idVal.(string)
		realNameStr, realNameOK := realNameVal.(string)

		if !idOK || !realNameOK {
			logx.Errorf("[TokenValidate] 类型转换失败 - idType: %T, realNameType: %T", idVal, realNameVal)
			err := errorx.New(errorx.SystemError, "TYPE_CONVERSION", "用户信息类型错误")
			httpx.WriteJson(w, http.StatusUnauthorized, GatewayErrorResult{
				Status:  http.StatusUnauthorized,
				Message: err.Message,
			})
			return
		}

		// 创建用户上下文信息对象
		ctxUserInfo := &types.UserInfo{
			ID:       idStr,
			Username: username,
			RealName: realNameStr,
		}

		// 记录最终要添加到上下文的用户信息
		logx.Infof("[TokenValidate] 添加到上下文的用户信息: %+v", ctxUserInfo)

		// 将用户信息添加到请求上下文中
		ctx := context.WithValue(r.Context(), types.UserContextKey, ctxUserInfo)

		// 使用更新后的上下文创建新的请求对象
		r = r.WithContext(ctx)

		// 记录验证成功
		logx.Infof("[TokenValidate] 令牌验证成功 - 用户: '%s'", username)

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

type GatewayErrorResult struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
