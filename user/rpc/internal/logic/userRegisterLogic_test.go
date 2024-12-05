package logic

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"flag"
	"shorterurl/user/rpc/internal/types/errorx"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"

	"shorterurl/user/rpc/internal/config"
	"shorterurl/user/rpc/internal/svc"
	__ "shorterurl/user/rpc/pb"
)

// generateTestUsername 生成随机测试用户名
func generateTestUsername() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "test_" + hex.EncodeToString([]byte(time.Now().String()))
	}
	return "test_" + hex.EncodeToString(bytes)
}

var (
	once   sync.Once
	svcCtx *svc.ServiceContext
	ctx    context.Context
)

// setupTest 设置测试环境
func setupTest(t *testing.T) (*svc.ServiceContext, context.Context) {
	once.Do(func() {
		configFile := flag.String("f", "../../etc/user.yaml", "配置文件路径")
		flag.Parse()

		var c config.Config
		conf.MustLoad(*configFile, &c)

		svcCtx = svc.NewServiceContext(c)
		logx.Disable()
		ctx = context.Background()
	})

	return svcCtx, ctx
}

// TestUserRegister 测试用户注册逻辑
func TestUserRegister(t *testing.T) {
	svcCtx, ctx := setupTest(t)
	logic := NewUserRegisterLogic(ctx, svcCtx)

	t.Run("成功注册新用户", func(t *testing.T) {
		req := &__.RegisterRequest{
			Username: generateTestUsername(),
			Password: "password123",
			RealName: "Test User",
			Phone:    "13800138000",
			Mail:     "test@example.com",
		}

		resp, err := logic.UserRegister(req)
		require.NoError(t, err, "注册应该成功")
		require.NotNil(t, resp, "响应不应为空")

		assert.Equal(t, req.Username, resp.Username)
		assert.NotEmpty(t, resp.CreateTime)
		assert.Equal(t, "注册成功", resp.Message)

		// 验证布隆过滤器
		exists, err := svcCtx.BloomFilters.UserExists(ctx, req.Username)
		require.NoError(t, err, "布隆过滤器检查失败")
		assert.True(t, exists, "用户应该存在于布隆过滤器中")
	})

	t.Run("重复注册用户名", func(t *testing.T) {
		username := generateTestUsername()
		firstReq := &__.RegisterRequest{
			Username: username,
			Password: "password123",
			RealName: "Test User",
			Phone:    "13800138000",
			Mail:     "test@example.com",
		}

		// 第一次注册
		resp1, err := logic.UserRegister(firstReq)
		require.NoError(t, err, "第一次注册应该成功")
		require.NotNil(t, resp1, "第一次注册响应不应为空")

		// 尝试重复注册
		resp2, err := logic.UserRegister(firstReq)
		require.Error(t, err, "重复注册应该失败")
		assert.Nil(t, resp2, "重复注册响应应该为空")

		// 检查错误类型
		appErr, ok := err.(*errorx.AppError)
		assert.True(t, ok, "应该返回 AppError")
		assert.Equal(t, errorx.ClientError, appErr.Type, "应该是客户端错误")
		assert.Equal(t, errorx.ErrUserNameExists, appErr.Code, "应该是用户名已存在错误")

		// 验证布隆过滤器
		exists, err := svcCtx.BloomFilters.UserExists(ctx, username)
		require.NoError(t, err, "布隆过滤器检查失败")
		assert.True(t, exists, "用户应该存在于布隆过滤器中")
	})

	t.Run("并发注册相同用户名", func(t *testing.T) {
		username := generateTestUsername()
		concurrency := 10
		type result struct {
			resp *__.RegisterResponse
			err  error
		}
		results := make(chan result, concurrency)

		// 并发注册
		for i := 0; i < concurrency; i++ {
			go func() {
				req := &__.RegisterRequest{
					Username: username,
					Password: "password123",
					RealName: "Test User",
					Phone:    "13800138000",
					Mail:     "test@example.com",
				}
				resp, err := logic.UserRegister(req)
				results <- result{resp, err}
			}()
		}

		// 统计结果
		successCount := 0
		clientErrorCount := 0
		systemErrorCount := 0

		for i := 0; i < concurrency; i++ {
			r := <-results
			if r.err == nil {
				successCount++
				continue
			}

			var appErr *errorx.AppError
			if errors.As(r.err, &appErr) {
				switch appErr.Type {
				case errorx.ClientError:
					clientErrorCount++
				case errorx.SystemError:
					systemErrorCount++
				}
			}
		}

		// 验证结果
		assert.Equal(t, 1, successCount, "应该只有一次注册成功")
		assert.Equal(t, concurrency-1, clientErrorCount, "其他应该都是客户端错误")
		assert.Equal(t, 0, systemErrorCount, "不应该有系统错误")

		// 验证布隆过滤器
		exists, err := svcCtx.BloomFilters.UserExists(ctx, username)
		require.NoError(t, err, "布隆过滤器检查失败")
		assert.True(t, exists, "用户应该存在于布隆过滤器中")
	})

	t.Run("布隆过滤器已存在用户名", func(t *testing.T) {
		username := generateTestUsername()

		// 先添加到布隆过滤器
		err := svcCtx.BloomFilters.AddUser(ctx, username)
		require.NoError(t, err, "添加到布隆过滤器应该成功")

		// 尝试注册
		req := &__.RegisterRequest{
			Username: username,
			Password: "password123",
			RealName: "Test User",
			Phone:    "13800138000",
			Mail:     "test@example.com",
		}

		resp, err := logic.UserRegister(req)
		require.Error(t, err, "注册应该失败")
		assert.Nil(t, resp, "响应应该为空")

		// 检查错误类型
		appErr, ok := err.(*errorx.AppError)
		assert.True(t, ok, "应该返回 AppError")
		assert.Equal(t, errorx.ClientError, appErr.Type, "应该是客户端错误")
		assert.Equal(t, errorx.ErrUserNameExists, appErr.Code, "应该是用户名已存在错误")
	})
}
