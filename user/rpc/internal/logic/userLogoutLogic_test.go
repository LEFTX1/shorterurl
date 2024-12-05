package logic

import (
	"errors"
	"shorterurl/user/rpc/internal/constant"
	"shorterurl/user/rpc/internal/types/errorx"
	__ "shorterurl/user/rpc/pb"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserLogout(t *testing.T) {
	svcCtx, ctx := setupTest(t)
	logic := NewUserLogoutLogic(ctx, svcCtx)
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

	t.Run("成功退出登录", func(t *testing.T) {
		req := &__.LogoutRequest{
			Username: username,
			Token:    loginResp.Token,
		}

		resp, err := logic.UserLogout(req)
		require.NoError(t, err, "退出登录应该成功")
		require.NotNil(t, resp, "响应不应为空")

		// 验证Redis中的数据已被删除
		exists, err := svcCtx.Redis.ExistsCtx(ctx, constant.UserLoginKey+username)
		require.NoError(t, err, "检查Redis数据失败")
		assert.False(t, exists, "Redis中的登录信息应该已被删除")
	})

	t.Run("用户未登录", func(t *testing.T) {
		req := &__.LogoutRequest{
			Username: "nonexistent_user",
			Token:    "invalid_token",
		}

		resp, err := logic.UserLogout(req)
		require.Error(t, err, "应该返回错误")
		assert.Nil(t, resp, "响应应该为空")

		// 验证错误类型
		var appErr *errorx.AppError
		ok := errors.As(err, &appErr)
		assert.True(t, ok, "应该返回 AppError")
		assert.Equal(t, errorx.ClientError, appErr.Type)
		assert.Equal(t, errorx.ErrUserNotFound, appErr.Code)
	})

	t.Run("Token无效", func(t *testing.T) {
		// 先登录用户
		_, err := loginLogic.UserLogin(&__.LoginRequest{
			Username: username,
			Password: registerReq.Password,
		})
		require.NoError(t, err, "登录失败")

		// 使用错误的token尝试退出
		req := &__.LogoutRequest{
			Username: username,
			Token:    "invalid_token",
		}

		resp, err := logic.UserLogout(req)
		require.Error(t, err, "应该返回错误")
		assert.Nil(t, resp, "响应应该为空")

		// 验证错误类型
		var appErr *errorx.AppError
		ok := errors.As(err, &appErr)
		assert.True(t, ok, "应该返回 AppError")
		assert.Equal(t, errorx.ClientError, appErr.Type)
		assert.Equal(t, errorx.ErrUserNotFound, appErr.Code)
	})
}
