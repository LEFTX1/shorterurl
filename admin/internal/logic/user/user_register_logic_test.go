// internal/logic/user/user_register_logic_test.go
package user

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"go-zero-shorterurl/admin/internal/config"
	"go-zero-shorterurl/admin/internal/svc"
	"go-zero-shorterurl/admin/internal/types"
	"go-zero-shorterurl/admin/internal/types/errorx"
	"testing"
	"time"

	"github.com/zeromicro/go-zero/core/conf"
)

// generateTestUsername 生成测试用户名
func generateTestUsername() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "test_" + hex.EncodeToString([]byte(time.Now().String()))
	}
	return "test_" + hex.EncodeToString(bytes)
}

func TestUserRegister(t *testing.T) {
	// 1. 加载配置
	configFile := flag.String("f", "../../../etc/admin-api.yaml", "the config file")
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 2. 创建服务上下文
	svcCtx, err := svc.NewServiceContext(c)
	if err != nil {
		t.Fatalf("create service context failed: %v", err)
	}
	defer func(svcCtx *svc.ServiceContext) {
		err = svcCtx.Close()
		if err != nil {
			t.Fatalf("close service context failed: %v", err)
		}
	}(svcCtx)

	// 3. 创建测试逻辑对象
	l := NewUserRegisterLogic(context.Background(), svcCtx)

	// 4. 测试正常注册
	t.Run("成功注册新用户", func(t *testing.T) {
		req := &types.UserRegisterReq{
			Username: generateTestUsername(),
			Password: "password123",
			RealName: "Test User",
			Phone:    "13800138000",
			Mail:     "test@example.com",
		}

		err = l.UserRegister(req)
		if err != nil {
			t.Errorf("UserRegister() error = %v, want nil", err)
		}
	})

	// 5. 测试重复注册
	t.Run("重复注册用户名", func(t *testing.T) {
		// 先注册一个用户
		username := generateTestUsername()
		firstReq := &types.UserRegisterReq{
			Username: username,
			Password: "password123",
		}
		err = l.UserRegister(firstReq)
		if err != nil {
			t.Fatalf("failed to register first user: %v", err)
		}

		// 尝试重复注册
		duplicateReq := &types.UserRegisterReq{
			Username: username,
			Password: "password123",
			RealName: "Test User",
			Phone:    "13800138000",
			Mail:     "test@example.com",
		}

		err = l.UserRegister(duplicateReq)
		if err == nil {
			t.Error("UserRegister() error = nil, want UserNameExistError")
		}
		if err.Error() != errorx.NewUserError(errorx.UserNameExistError).Error() {
			t.Errorf("UserRegister() error = %v, want %v", err, errorx.NewUserError(errorx.UserNameExistError))
		}
	})

	// 6. 测试并发注册
	t.Run("并发注册相同用户名", func(t *testing.T) {
		username := generateTestUsername()
		concurrency := 10
		done := make(chan error, concurrency)

		// 并发注册同一个用户名
		for i := 0; i < concurrency; i++ {
			go func() {
				req := &types.UserRegisterReq{
					Username: username,
					Password: "password123",
					RealName: "Test User",
				}
				done <- l.UserRegister(req)
			}()
		}

		// 统计成功次数
		successCount := 0
		for i := 0; i < concurrency; i++ {
			err = <-done
			if err == nil {
				successCount++
			}
		}

		// 应该只有一个注册成功
		if successCount != 1 {
			t.Errorf("expected 1 successful registration, got %d", successCount)
		}
	})
}
