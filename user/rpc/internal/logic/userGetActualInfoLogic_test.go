package logic

import (
	"errors"
	"shorterurl/user/rpc/internal/types/errorx"
	__ "shorterurl/user/rpc/pb"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserGetActualInfo(t *testing.T) {
	svcCtx, ctx := setupTest(t)
	logic := NewUserGetActualInfoLogic(ctx, svcCtx)
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

	// 注册用户
	_, err := registerLogic.UserRegister(registerReq)
	require.NoError(t, err, "注册用户失败")

	t.Run("成功获取用户信息", func(t *testing.T) {
		req := &__.CheckUsernameRequest{
			Username: username,
		}

		resp, err := logic.UserGetActualInfo(req)
		require.NoError(t, err, "获取用户信息应该成功")
		require.NotNil(t, resp, "响应不应为空")

		// 验证返回的用户信息（无脱敏）
		assert.Equal(t, username, resp.Username)
		assert.Equal(t, registerReq.RealName, resp.RealName)
		assert.Equal(t, registerReq.Phone, resp.Phone)
		assert.Equal(t, registerReq.Mail, resp.Mail)
		assert.NotEmpty(t, resp.CreateTime)
		assert.NotEmpty(t, resp.UpdateTime)
	})

	t.Run("用户不存在", func(t *testing.T) {
		req := &__.CheckUsernameRequest{
			Username: "nonexistent_user",
		}

		resp, err := logic.UserGetActualInfo(req)
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
