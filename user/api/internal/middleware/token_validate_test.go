package middleware

import (
	"encoding/json"     // 导入 JSON 编码/解码包
	"flag"              // 导入命令行标志解析包
	"net/http"          // 导入 HTTP 客户端和服务器包
	"net/http/httptest" // 导入 HTTP 测试包
	"testing"           // 导入测试包

	"github.com/stretchr/testify/assert"             // 导入断言包
	"github.com/zeromicro/go-zero/core/conf"         // 导入配置加载包
	"github.com/zeromicro/go-zero/core/stores/redis" // 导入 Redis 客户端包
	"github.com/zeromicro/go-zero/rest"              // 导入 REST 服务器包

	"shorterurl/user/api/internal/config" // 导入内部配置包
	"shorterurl/user/api/internal/svc"    // 导入内部服务上下文包
	"shorterurl/user/api/internal/types"  // 导入内部类型包
)

// 配置文件路径
var configFile = flag.String("f", "../../etc/gateway.yaml", "配置文件路径")

// setupTest 初始化测试环境
// 返回值1: TokenValidateMiddleware 实例
// 返回值2: Redis 客户端实例
func setupTest() (*TokenValidateMiddleware, *redis.Redis) {
	flag.Parse() // 解析命令行标志

	// 加载配置文件
	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 创建服务器实例
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop() // 确保服务器在函数结束时停止

	// 创建服务上下文
	svcCtx := svc.NewServiceContext(c)

	// 准备测试用户数据
	userInfo := map[string]string{
		"id":       "123",
		"realName": "测试用户",
	}

	// 将用户信息序列化并存储到 Redis
	userInfoBytes, _ := json.Marshal(userInfo)
	err := svcCtx.Redis.Hset("short-link:login:testuser", "valid-token", string(userInfoBytes))
	if err != nil {
		panic("设置用户信息失败: " + err.Error())
	}

	return NewTokenValidateMiddleware(&c.Auth, svcCtx.Redis), svcCtx.Redis
}

// TestTokenValidateMiddleware 中间件测试主函数
func TestTokenValidateMiddleware(t *testing.T) {
	// 测试白名单路径
	t.Run("TestWhiteListPath", func(t *testing.T) {
		middleware, _ := setupTest() // 初始化测试环境

		// 创建登录路径的测试请求
		req := httptest.NewRequest(http.MethodGet, "/api/short-link/admin/v1/user/login", nil)
		w := httptest.NewRecorder() // 创建响应记录器

		var called bool
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true // 标记处理程序被调用
		})

		middleware.Handle(next).ServeHTTP(w, req) // 调用中间件处理请求

		assert.True(t, called, "白名单路径应该被正常处理")
		assert.Equal(t, http.StatusOK, w.Code) // 断言响应状态码为 200
	})

	// 测试缺少凭证的情况
	t.Run("TestMissingCredentials", func(t *testing.T) {
		middleware, _ := setupTest() // 初始化测试环境

		// 创建一个受保护路径的请求
		req := httptest.NewRequest(http.MethodGet, "/api/protected", nil)
		w := httptest.NewRecorder() // 创建响应记录器

		var called bool
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true // 标记处理程序被调用
		})

		middleware.Handle(next).ServeHTTP(w, req) // 调用中间件处理请求

		assert.False(t, called, "缺少凭证时不应该调用处理函数")
		assert.Equal(t, http.StatusUnauthorized, w.Code) // 断言响应状态码为 401

		var response types.GatewayErrorResult
		err := json.NewDecoder(w.Body).Decode(&response) // 解码响应体
		assert.NoError(t, err)                           // 断言解码无错误
		assert.Equal(t, "用户名或令牌不能为空", response.Message)  // 断言错误消息
	})

	// 测试无效的令牌
	t.Run("TestInvalidToken", func(t *testing.T) {
		middleware, _ := setupTest() // 初始化测试环境

		// 创建带有无效令牌的请求
		req := httptest.NewRequest(http.MethodGet, "/api/protected", nil)
		req.Header.Set("username", "testuser")   // 设置用户名头
		req.Header.Set("token", "invalid-token") // 设置无效令牌头
		w := httptest.NewRecorder()              // 创建响应记录器

		var called bool
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true // 标记处理程序被调用
		})

		middleware.Handle(next).ServeHTTP(w, req) // 调用中间件处理请求

		assert.False(t, called, "无效令牌时不应该调用处理函数")
		assert.Equal(t, http.StatusUnauthorized, w.Code) // 断言响应状态码为 401

		var response types.GatewayErrorResult
		err := json.NewDecoder(w.Body).Decode(&response) // 解码响应体
		assert.NoError(t, err)                           // 断言解码无错误
		assert.Equal(t, "无效的令牌", response.Message)       // 断言错误消息
	})

	// 测试有效的令牌
	t.Run("TestValidToken", func(t *testing.T) {
		middleware, _ := setupTest() // 初始化测试环境

		// 创建带有有效令牌的请求
		req := httptest.NewRequest(http.MethodGet, "/api/protected", nil)
		req.Header.Set("username", "testuser") // 设置用户名头
		req.Header.Set("token", "valid-token") // 设置有效令牌头
		w := httptest.NewRecorder()            // 创建响应记录器

		var called bool
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true // 标记处理程序被调用
			// 验证上下文中的用户信息
			userInfo, ok := types.GetUserFromCtx(r.Context())
			assert.True(t, ok, "上下文中应该包含用户信息")
			assert.Equal(t, "123", userInfo.ID)            // 断言用户 ID
			assert.Equal(t, "testuser", userInfo.Username) // 断言用户名
			assert.Equal(t, "测试用户", userInfo.RealName)     // 断言真实姓名
		})

		middleware.Handle(next).ServeHTTP(w, req) // 调用中间件处理请求

		assert.True(t, called, "有效令牌时应该调用处理函数")
		assert.Equal(t, http.StatusOK, w.Code) // 断言响应状态码为 200
	})
}
