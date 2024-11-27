// internal/logic/user/user_info_logic_test.go
package user

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"shorterurl/admin/internal/dal/model"
	"shorterurl/admin/internal/types"
	"shorterurl/admin/internal/types/errorx"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestUserInfoLogic(t *testing.T) {
	svcCtx, ctx := setupTest(t)

	// 1. 先注册一个用户
	registerLogic := NewUserRegisterLogic(ctx, svcCtx)
	testUsername := generateTestUsername()
	registerReq := &types.UserRegisterReq{
		Username: testUsername,
		Password: "password123",
		RealName: "Test User",
		Phone:    "13800138000",
		Mail:     "test@example.com",
	}

	// 注册用户
	registerResp, err := registerLogic.UserRegister(registerReq)
	require.NoError(t, err, "注册用户失败")
	require.NotNil(t, registerResp, "注册响应不应为空")

	// 等待数据写入
	time.Sleep(time.Second)

	// 2. 测试用户信息查询
	infoLogic := NewUserInfoLogic(ctx, svcCtx)

	t.Run("查询空用户名", func(t *testing.T) {
		resp, err := infoLogic.UserInfo(&types.UserUsernameReq{Username: ""})
		require.Error(t, err, "空用户名应该返回错误")
		require.Nil(t, resp, "响应应该为空")

		var userErr *errorx.UserError
		if assert.True(t, errors.As(err, &userErr), "应该是用户错误") {
			assert.Equal(t, "USERNAME_EMPTY", userErr.Code)
		}
	})

	t.Run("查询不存在的用户", func(t *testing.T) {
		// 先确认用户确实不存在
		var nonExistUser model.TUser
		err := svcCtx.DB.Where("username = ? AND del_flag = ?", "nonexistent_user", false).First(&nonExistUser).Error
		t.Logf("不存在用户验证: err=%v", err)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound, "应该返回记录未找到错误")

		resp, err := infoLogic.UserInfo(&types.UserUsernameReq{
			Username: "nonexistent_user",
		})
		require.Error(t, err, "不存在的用户应该返回错误")
		require.Nil(t, resp, "响应应该为空")

		var userErr *errorx.UserError
		if assert.True(t, errors.As(err, &userErr), "应该是用户错误") {
			assert.Equal(t, "USER_NOT_EXIST", userErr.Code)
		}
	})

	t.Run("查询已注册用户", func(t *testing.T) {
		// 先确认用户确实存在
		var existingUser model.TUser
		err := svcCtx.DB.Where("username = ? AND del_flag = ?", testUsername, false).First(&existingUser).Error
		require.NoError(t, err, "查询已注册用户失败")
		t.Logf("已注册用户数据: %+v", existingUser)

		resp, err := infoLogic.UserInfo(&types.UserUsernameReq{
			Username: testUsername,
		})
		require.NoError(t, err, "查询已注册用户不应该返回错误")
		require.NotNil(t, resp, "响应不应为空")

		// 验证返回的用户信息
		assert.Equal(t, testUsername, resp.Username)
		assert.Equal(t, registerReq.RealName, resp.RealName)

		// **验证脱敏后的手机号和邮箱**
		expectedMaskedPhone := "138****8000" // 手机号脱敏后的格式
		assert.Equal(t, expectedMaskedPhone, resp.Phone, "手机号应该是脱敏后的数据")

		expectedMaskedEmail := "t****@example.com" // 邮箱脱敏后的格式
		assert.Equal(t, expectedMaskedEmail, resp.Mail, "邮箱应该是脱敏后的数据")

		assert.NotEmpty(t, resp.CreateTime)
		assert.NotEmpty(t, resp.UpdateTime)
	})

}

func captureOutput(f func()) string {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	f()

	w.Close()
	os.Stdout = oldStdout
	return <-outC
}

func TestSQLGeneration(t *testing.T) {
	// 1. 使用 setupTest 获取服务上下文
	svcCtx, ctx := setupTest(t)

	// 2. 准备测试数据
	username := "test_user"

	// 3. 测试原生 GORM 方式
	t.Run("GORM Native Query", func(t *testing.T) {
		var user model.TUser
		svcCtx.DB.
			Where("username = ? AND del_flag = ?", username, false).
			First(&user)
		t.Logf("user: %+v", user)
		t.Logf("GORM SQL (ending): %v", svcCtx.Sharding.LastQuery())

	})

	// 4. 测试 Gen 使用 Select
	t.Run("Gen Query With Select", func(t *testing.T) {
		output := captureOutput(func() {
			// 执行查询以触发 Debug 输出
			_, err := svcCtx.Query.TUser.Debug().WithContext(ctx).
				Select(svcCtx.Query.TUser.Username, svcCtx.Query.TUser.DelFlag).
				Where(svcCtx.Query.TUser.Username.Eq(username)).
				Where(svcCtx.Query.TUser.DelFlag.Is(false)).
				First()

			if err != nil {
				t.Logf("Expected error: %v", err)
			}
		})
		t.Logf("Gen SQL: %v", svcCtx.Sharding.LastQuery())
		// 同时在终端显示
		fmt.Print(output)
	})

	// 5. 对比：不使用 Select 的 Gen 查询
	t.Run("Gen Query Without Select", func(t *testing.T) {
		output := captureOutput(func() {
			// 执行查询以触发 Debug 输出
			_, err := svcCtx.Query.TUser.Debug().WithContext(ctx).
				Where(svcCtx.Query.TUser.Username.Eq(username)).
				Where(svcCtx.Query.TUser.DelFlag.Is(false)).
				First()

			if err != nil {
				t.Logf("Expected error: %v", err)
			}
		})

		// 显示捕获的输出
		t.Log("Terminal Output:", output)
		t.Logf("Gen SQL: %v", svcCtx.Sharding.LastQuery())
		// 同时在终端显示
		fmt.Print(output)
	})

}

// 辅助函数：格式化 SQL 和参数
func formatSQL(sql string, vars []interface{}) string {
	for _, v := range vars {
		switch v := v.(type) {
		case string:
			sql = strings.Replace(sql, "?", "'"+v+"'", 1)
		case bool:
			sql = strings.Replace(sql, "?", fmt.Sprintf("%v", v), 1)
		default:
			sql = strings.Replace(sql, "?", fmt.Sprintf("%v", v), 1)
		}
	}
	return sql
}
