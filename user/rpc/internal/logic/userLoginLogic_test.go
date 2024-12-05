package logic

import (
	"shorterurl/user/rpc/internal/constant"
	"shorterurl/user/rpc/internal/dal/query"
	"shorterurl/user/rpc/internal/types/errorx"
	__ "shorterurl/user/rpc/pb"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserLogin(t *testing.T) {
	svcCtx, ctx := setupTest(t)
	logic := NewUserLoginLogic(ctx, svcCtx)
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

	// 清理测试数据
	defer func() {
		q := query.Use(svcCtx.DB)
		_, _ = q.TUser.WithContext(ctx).Where(q.TUser.Username.Eq(username)).Delete()
		_, _ = svcCtx.Redis.DelCtx(ctx, constant.USER_LOGIN_KEY+username)
	}()

	t.Run("用户不存在", func(t *testing.T) {
		req := &__.LoginRequest{
			Username: "nonexistent_user",
			Password: "password123",
		}

		resp, err := logic.UserLogin(req)
		require.Error(t, err, "应该返回错误")
		assert.Nil(t, resp, "响应应该为空")

		// 验证错误类型
		appErr, ok := err.(*errorx.AppError)
		assert.True(t, ok, "应该返回 AppError")
		assert.Equal(t, errorx.ClientError, appErr.Type)
		assert.Equal(t, errorx.ErrUserNotFound, appErr.Code)
	})

	t.Run("成功登录", func(t *testing.T) {
		// 先注册用户
		_, err := registerLogic.UserRegister(registerReq)
		require.NoError(t, err, "注册用户失败")

		// 登录请求
		loginReq := &__.LoginRequest{
			Username: username,
			Password: registerReq.Password,
		}

		// 第一次登录
		resp1, err := logic.UserLogin(loginReq)
		require.NoError(t, err, "第一次登录应该成功")
		require.NotNil(t, resp1, "响应不应为空")
		assert.Equal(t, username, resp1.Username)
		assert.Equal(t, registerReq.RealName, resp1.RealName)
		assert.NotEmpty(t, resp1.Token)
		assert.NotEmpty(t, resp1.CreateTime)

		// 验证Redis中的数据
		val, err := svcCtx.Redis.HgetallCtx(ctx, constant.USER_LOGIN_KEY+username)
		require.NoError(t, err, "获取Redis数据失败")
		assert.NotEmpty(t, val, "Redis中应该有数据")
		assert.Equal(t, 1, len(val), "应该只有一个token")
		assert.Equal(t, username, val[resp1.Token], "Redis中的用户名应该匹配")

		// 第二次登录（使用相同token）
		resp2, err := logic.UserLogin(loginReq)
		require.NoError(t, err, "第二次登录应该成功")
		require.NotNil(t, resp2, "响应不应为空")
		assert.Equal(t, resp1.Token, resp2.Token, "应该返回相同的token")
	})

	t.Run("密码错误", func(t *testing.T) {
		req := &__.LoginRequest{
			Username: username,
			Password: "wrong_password",
		}

		resp, err := logic.UserLogin(req)
		require.Error(t, err, "密码错误应该返回错误")
		assert.Nil(t, resp, "响应应该为空")

		// 验证错误类型
		appErr, ok := err.(*errorx.AppError)
		assert.True(t, ok, "应该返回 AppError")
		assert.Equal(t, errorx.ClientError, appErr.Type)
		assert.Equal(t, errorx.ErrUserNotFound, appErr.Code)
	})
}
