package logic

import (
	"shorterurl/user/rpc/internal/types/errorx"
	__ "shorterurl/user/rpc/pb"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserUpdate(t *testing.T) {
	svcCtx, ctx := setupTest(t)
	logic := NewUserUpdateLogic(ctx, svcCtx)
	registerLogic := NewUserRegisterLogic(ctx, svcCtx)
	getInfoLogic := NewUserGetActualInfoLogic(ctx, svcCtx)

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

	t.Run("成功更新用户信息", func(t *testing.T) {
		req := &__.UpdateRequest{
			Username: username,
			Password: "newpassword123",
			RealName: "New Test User",
			Phone:    "13900139000",
			Mail:     "newtest@example.com",
		}

		resp, err := logic.UserUpdate(req)
		require.NoError(t, err, "更新用户信息应该成功")
		require.NotNil(t, resp, "响应不应为空")
		assert.True(t, resp.Success, "更新应该成功")
		assert.Equal(t, "更新成功", resp.Message)

		// 验证更新后的信息
		info, err := getInfoLogic.UserGetActualInfo(&__.CheckUsernameRequest{Username: username})
		require.NoError(t, err, "获取用户信息应该成功")
		assert.Equal(t, req.RealName, info.RealName)
		assert.Equal(t, req.Phone, info.Phone)
		assert.Equal(t, req.Mail, info.Mail)
	})

	t.Run("部分更新用户信息", func(t *testing.T) {
		req := &__.UpdateRequest{
			Username: username,
			RealName: "Partial Update User",
		}

		resp, err := logic.UserUpdate(req)
		require.NoError(t, err, "更新用户信息应该成功")
		require.NotNil(t, resp, "响应不应为空")
		assert.True(t, resp.Success, "更新应该成功")

		// 验证更新后的信息
		info, err := getInfoLogic.UserGetActualInfo(&__.CheckUsernameRequest{Username: username})
		require.NoError(t, err, "获取用户信息应该成功")
		assert.Equal(t, req.RealName, info.RealName)
		// 其他字段应该保持不变
		assert.Equal(t, "13900139000", info.Phone)
		assert.Equal(t, "newtest@example.com", info.Mail)
	})

	t.Run("用户不存在", func(t *testing.T) {
		req := &__.UpdateRequest{
			Username: "nonexistent_user",
			RealName: "New Name",
		}

		resp, err := logic.UserUpdate(req)
		require.Error(t, err, "应该返回错误")
		assert.Nil(t, resp, "响应应该为空")

		// 验证错误类型
		appErr, ok := err.(*errorx.AppError)
		assert.True(t, ok, "应该返回 AppError")
		assert.Equal(t, errorx.ClientError, appErr.Type)
		assert.Equal(t, errorx.ErrUserNotFound, appErr.Code)
	})
}
