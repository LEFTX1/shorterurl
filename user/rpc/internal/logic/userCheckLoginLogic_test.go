package logic

import (
	"shorterurl/user/rpc/internal/constant"
	__ "shorterurl/user/rpc/pb"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserCheckLogin(t *testing.T) {
	svcCtx, ctx := setupTest(t)
	logic := NewUserCheckLoginLogic(ctx, svcCtx)
	loginLogic := NewUserLoginLogic(ctx, svcCtx)
	registerLogic := NewUserRegisterLogic(ctx, svcCtx)

	// 准备测试数据：用户名和注册请求
	username := generateTestUsername()
	registerReq := &__.RegisterRequest{
		Username: username,
		Password: "password123",
		RealName: "Test User",
		Phone:    "13800138000",
		Mail:     "test@example.com",
	}

	// 注册并登录用户
	_, err := registerLogic.UserRegister(registerReq)
	require.NoError(t, err, "注册用户失败")

	loginResp, err := loginLogic.UserLogin(&__.LoginRequest{
		Username: username,
		Password: registerReq.Password,
	})
	require.NoError(t, err, "登录失败")
	require.NotNil(t, loginResp, "登录响应不应为空")

	t.Run("成功检查登录状态", func(t *testing.T) {
		req := &__.CheckUsernameRequest{
			Username: username,
		}

		resp, err := logic.UserCheckLogin(req)
		require.NoError(t, err, "检查登录状态应该成功")
		require.NotNil(t, resp, "响应不应为空")
		assert.True(t, resp.Success, "用户应该处于登录状态")
		assert.Equal(t, "用户已登录", resp.Message)

		// 验证Redis中的数据仍然存在
		exists, err := svcCtx.Redis.ExistsCtx(ctx, constant.USER_LOGIN_KEY+username)
		require.NoError(t, err, "检查Redis数据失败")
		assert.True(t, exists, "Redis中的登录信息应该存在")
	})

	t.Run("用户未登录", func(t *testing.T) {
		req := &__.CheckUsernameRequest{
			Username: "nonexistent_user",
		}

		resp, err := logic.UserCheckLogin(req)
		require.NoError(t, err, "应该返回成功但未登录的响应")
		require.NotNil(t, resp, "响应不应为空")
		assert.False(t, resp.Success, "用户应该未登录")
		assert.Equal(t, "用户未登录", resp.Message)
	})
}
