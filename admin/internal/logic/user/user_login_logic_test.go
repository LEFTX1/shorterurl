// internal/logic/user/user_login_logic_test.go
package user

import (
	"context"
	"shorterurl/admin/internal/dal/model"
	"shorterurl/admin/internal/svc"
	"shorterurl/admin/internal/types"
	"shorterurl/admin/internal/types/errorx"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserLoginLogic(t *testing.T) {
	svcCtx, ctx := setupTest(t)
	logic := NewUserLoginLogic(ctx, svcCtx)

	t.Run("成功登录", func(t *testing.T) {
		username := generateTestUsername()
		password := "password123"
		createTestUser(t, svcCtx, username, password)
		defer cleanupTestUser(t, svcCtx, username)

		resp, err := logic.UserLogin(&types.UserLoginReq{
			Username: username,
			Password: password,
		})

		require.NoError(t, err, "登录应该成功")
		require.NotNil(t, resp, "响应不应为空")
		assert.NotEmpty(t, resp.Token, "token不应为空")

		loginKey := "login:" + username
		exists, err := svcCtx.Redis.ExistsCtx(ctx, loginKey)
		require.NoError(t, err, "检查Redis key失败")
		assert.True(t, exists, "登录信息应该存在于Redis中")

		ttl, err := svcCtx.Redis.TtlCtx(ctx, loginKey)
		require.NoError(t, err, "获取过期时间失败")
		assert.True(t, ttl <= loginExpireTime && ttl > 0, "过期时间应该在合理范围内")
	})

	t.Run("用户名不存在", func(t *testing.T) {
		resp, err := logic.UserLogin(&types.UserLoginReq{
			Username: "nonexistent_user",
			Password: "password123",
		})

		assert.Nil(t, resp, "响应应该为空")
		assert.Error(t, err, "应该返回错误")

		userErr, ok := err.(*errorx.UserError)
		assert.True(t, ok, "应该是用户错误")
		assert.Equal(t, errorx.UserNotExistError, userErr.Code, "错误码应该是用户不存在")
	})

	t.Run("密码错误", func(t *testing.T) {
		username := generateTestUsername()
		password := "password123"
		createTestUser(t, svcCtx, username, password)
		defer cleanupTestUser(t, svcCtx, username)

		resp, err := logic.UserLogin(&types.UserLoginReq{
			Username: username,
			Password: "wrong_password",
		})

		assert.Nil(t, resp, "响应应该为空")
		assert.Error(t, err, "应该返回错误")

		userErr, ok := err.(*errorx.UserError)
		assert.True(t, ok, "应该是用户错误")
		assert.Equal(t, errorx.UserNotExistError, userErr.Code, "错误码应该是用户不存在")
	})

	t.Run("重复登录", func(t *testing.T) {
		username := generateTestUsername()
		password := "password123"
		createTestUser(t, svcCtx, username, password)
		defer cleanupTestUser(t, svcCtx, username)

		resp1, err := logic.UserLogin(&types.UserLoginReq{
			Username: username,
			Password: password,
		})
		require.NoError(t, err, "第一次登录应该成功")
		require.NotNil(t, resp1, "第一次登录响应不应为空")

		resp2, err := logic.UserLogin(&types.UserLoginReq{
			Username: username,
			Password: password,
		})
		require.NoError(t, err, "第二次登录应该成功")
		require.NotNil(t, resp2, "第二次登录响应不应为空")

		assert.Equal(t, resp1.Token, resp2.Token, "重复登录应该返回相同的token")

		loginKey := "login:" + username
		result, err := svcCtx.Redis.HgetallCtx(ctx, loginKey)
		require.NoError(t, err, "获取Redis数据失败")
		assert.Equal(t, 1, len(result), "Redis中应该只有一个token")
	})

	t.Run("并发登录", func(t *testing.T) {
		username := generateTestUsername()
		password := "password123"
		createTestUser(t, svcCtx, username, password)
		defer cleanupTestUser(t, svcCtx, username)

		concurrency := 10
		type loginResult struct {
			resp *types.UserLoginResp
			err  error
		}
		results := make(chan loginResult, concurrency)

		for i := 0; i < concurrency; i++ {
			go func() {
				resp, err := logic.UserLogin(&types.UserLoginReq{
					Username: username,
					Password: password,
				})
				results <- loginResult{resp, err}
			}()
		}

		var firstToken string
		successCount := 0
		for i := 0; i < concurrency; i++ {
			r := <-results
			if r.err == nil && r.resp != nil {
				successCount++
				if firstToken == "" {
					firstToken = r.resp.Token
				} else {
					assert.Equal(t, firstToken, r.resp.Token, "所有成功登录应该返回相同的token")
				}
			}
		}

		assert.Equal(t, concurrency, successCount, "所有登录请求都应该成功")

		loginKey := "login:" + username
		result, err := svcCtx.Redis.HgetallCtx(ctx, loginKey)
		require.NoError(t, err, "获取Redis数据失败")
		assert.Equal(t, 1, len(result), "Redis中应该只有一个token")
	})
}

// 辅助函数
func createTestUser(t *testing.T, svcCtx *svc.ServiceContext, username, password string) {
	user := &model.TUser{
		Username:   username,
		Password:   password,
		RealName:   "Test User",
		Phone:      "13800138000",
		Mail:       "test@example.com",
		DelFlag:    false,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	err := svcCtx.DB.Create(user).Error
	require.NoError(t, err, "创建测试用户失败")

	// 添加到布隆过滤器
	err = svcCtx.BloomFilters.AddUser(context.Background(), username)
	require.NoError(t, err, "添加用户到布隆过滤器失败")
}

func cleanupTestUser(t *testing.T, svcCtx *svc.ServiceContext, username string) {
	// 从数据库中删除用户
	err := svcCtx.DB.Unscoped().Where("username = ?", username).Delete(&model.TUser{}).Error
	require.NoError(t, err, "清理测试用户失败")

	// 删除Redis中的登录信息
	loginKey := "login:" + username
	_, err = svcCtx.Redis.DelCtx(context.Background(), loginKey)
	require.NoError(t, err, "清理Redis登录信息失败")
}
