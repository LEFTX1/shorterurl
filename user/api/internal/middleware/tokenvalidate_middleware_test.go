package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeromicro/go-zero/core/stores/redis"

	"shorterurl/user/api/internal/config"
)

func setupTest(t *testing.T) (*TokenValidateMiddleware, *redis.Redis) {
	c := config.Config{
		Auth: config.TokenValidateConfig{
			WhitePathList: []string{"/api/public", "/api/health"},
		},
		Redis: struct {
			RedisConf redis.RedisConf
		}{
			RedisConf: redis.RedisConf{
				Host: "localhost:6379",
				Type: "node",
			},
		},
	}

	redisClient := redis.MustNewRedis(c.Redis.RedisConf)
	middleware := NewTokenValidateMiddleware(&c.Auth, redisClient)

	return middleware, redisClient
}

func TestTokenValidateMiddleware(t *testing.T) {
	middleware, redisClient := setupTest(t)

	// 测试白名单路径
	t.Run("白名单路径应该直接放行", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/public/test", nil)
		w := httptest.NewRecorder()

		called := false
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
		})

		middleware.Handle(next).ServeHTTP(w, req)
		assert.True(t, called, "白名单路径应该调用下一个处理函数")
	})

	// 测试缺少凭证
	t.Run("缺少凭证应该返回未授权错误", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/protected", nil)
		w := httptest.NewRecorder()

		middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
			t.Error("不应该调用下一个处理函数")
		}).ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	// 测试无效的令牌
	t.Run("无效的令牌应该返回未授权错误", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/protected", nil)
		req.Header.Set("username", "testuser")
		req.Header.Set("token", "invalid_token")
		w := httptest.NewRecorder()

		middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
			t.Error("不应该调用下一个处理函数")
		}).ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	// 测试有效的令牌
	t.Run("有效的令牌应该成功通过", func(t *testing.T) {
		// 在 Redis 中设置测试数据
		userInfo := `{"id":"123","realName":"Test User"}`
		err := redisClient.Hset("user:login:testuser", "valid_token", userInfo)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/api/protected", nil)
		req.Header.Set("username", "testuser")
		req.Header.Set("token", "valid_token")
		w := httptest.NewRecorder()

		called := false
		middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
			called = true
		}).ServeHTTP(w, req)

		assert.True(t, called, "应该调用下一个处理函数")
	})
}
