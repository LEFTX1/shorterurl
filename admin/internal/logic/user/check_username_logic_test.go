package user

import (
	"shorterurl/admin/internal/types"
	"shorterurl/admin/internal/types/errorx"
	"testing"
)

func TestCheckUsernameLogic_CheckUsername(t *testing.T) {
	// 1 获取测试上下文
	svcCtx, ctx := setupTest(t)
	// 2 调用注册逻辑函数 其中用户名是随机生成
	// 2.1 生成用户名
	testUsername := generateTestUsername()

	// 2.2 调用注册逻辑函数
	req := &types.UserRegisterReq{
		Username: testUsername,
		Password: "password123",
		RealName: "Test User",
		Phone:    "13800138000",
		Mail:     "asd@qq.com",
	}
	registerLogic := NewUserRegisterLogic(ctx, svcCtx)
	_, err := registerLogic.UserRegister(req)
	if err != nil {
		t.Errorf("注册失败: %v", err)
	}

	// 3 利用生成的用户名调用checkusername函数 判断是否存在
	t.Run("CheckUsername", func(t *testing.T) {
		logic := NewCheckUsernameLogic(ctx, svcCtx)
		resp, err := logic.CheckUsername(&types.UserCheckUsernameReq{
			Username: testUsername,
		})
		if err != nil {
			t.Errorf("CheckUsername failed: %v", err)
		}
		if !resp {
			t.Errorf("expected username %s to exist, but CheckUsername returned false", testUsername)
		}
		t.Run("CheckUsername - Nonexistent Username", func(t *testing.T) {
			logic := NewCheckUsernameLogic(ctx, svcCtx)
			resp, err := logic.CheckUsername(&types.UserCheckUsernameReq{
				Username: generateTestUsername(),
			})
			if err != nil {
				t.Errorf("CheckUsername failed: %v", err)
			}
			if resp {
				t.Errorf("expected username nonexistent_user to not exist, but CheckUsername returned true")
			}
		})

		t.Run("CheckUsername - Empty Username", func(t *testing.T) {
			logic := NewCheckUsernameLogic(ctx, svcCtx)
			resp, err := logic.CheckUsername(&types.UserCheckUsernameReq{
				Username: "",
			})
			if err == nil {
				t.Errorf(errorx.NewUserError(errorx.UserNotExistError).Error())
			}
			if resp {
				t.Errorf("expected CheckUsername to return false for empty username, but got true")
			}
		})

	})

}
