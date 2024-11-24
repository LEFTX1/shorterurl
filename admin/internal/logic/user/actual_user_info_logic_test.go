package user

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go-zero-shorterurl/admin/internal/dal/model"
	"go-zero-shorterurl/admin/internal/types"
	"go-zero-shorterurl/admin/internal/types/errorx"
)

// TestActualUserInfoLogic 测试 ActualUserInfoLogic 接口
func TestActualUserInfoLogic(t *testing.T) {
	// 初始化测试上下文和服务
	svcCtx, ctx := setupTest(t)
	actualInfoLogic := NewActualUserInfoLogic(ctx, svcCtx)

	// 准备测试数据
	testUsername := "testuser"
	testPhone := "13800138000"
	testEmail := "test@example.com"

	// 插入测试用户
	err := svcCtx.DB.Create(&model.TUser{
		Username:   testUsername,
		Password:   "password123",
		RealName:   "Test User",
		Phone:      testPhone,
		Mail:       testEmail,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		DelFlag:    false,
	}).Error
	require.NoError(t, err, "初始化测试数据失败")

	t.Run("查询实际用户信息", func(t *testing.T) {
		// 构造请求
		req := &types.UserUsernameReq{
			Username: testUsername,
		}

		// 调用接口
		resp, err := actualInfoLogic.ActualUserInfo(req)
		require.NoError(t, err, "查询实际用户信息失败")
		require.NotNil(t, resp, "响应不应为空")

		// 验证返回值
		assert.Equal(t, testUsername, resp.Username)
		assert.Equal(t, "Test User", resp.RealName)
		assert.Equal(t, testPhone, resp.Phone, "手机号应该未脱敏")
		assert.Equal(t, testEmail, resp.Mail, "邮箱应该未脱敏")
		assert.NotEmpty(t, resp.CreateTime, "创建时间不应为空")
		assert.NotEmpty(t, resp.UpdateTime, "更新时间不应为空")
	})

	t.Run("查询不存在的用户", func(t *testing.T) {
		// 构造请求
		req := &types.UserUsernameReq{
			Username: "nonexistent_user",
		}

		// 调用接口
		resp, err := actualInfoLogic.ActualUserInfo(req)
		require.Error(t, err, "查询不存在的用户应该返回错误")
		require.Nil(t, resp, "响应应该为空")

		// 验证错误类型
		var userErr *errorx.UserError
		if assert.True(t, errors.As(err, &userErr), "应该是用户错误") {
			assert.Equal(t, "USER_NOT_EXIST", userErr.Code)
		}
	})

	t.Run("查询空用户名", func(t *testing.T) {
		// 构造请求
		req := &types.UserUsernameReq{
			Username: "",
		}

		// 调用接口
		resp, err := actualInfoLogic.ActualUserInfo(req)
		require.Error(t, err, "查询空用户名应该返回错误")
		require.Nil(t, resp, "响应应该为空")

		// 验证错误类型
		var userErr *errorx.UserError
		if assert.True(t, errors.As(err, &userErr), "应该是用户错误") {
			assert.Equal(t, "USERNAME_EMPTY", userErr.Code)
		}
	})
}
