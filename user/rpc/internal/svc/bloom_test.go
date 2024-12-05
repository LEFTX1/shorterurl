package svc

import (
	"context"
	"errors"
	"flag"
	"shorterurl/user/rpc/internal/config"
	"shorterurl/user/rpc/internal/types/errorx"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeromicro/go-zero/core/conf"
)

func TestBloomFilterManager_UserExists(t *testing.T) {
	configFile := flag.String("f", "../../etc/user.yaml", "配置文件路径")
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 初始化服务上下文
	svcCtx := NewServiceContext(c)

	// 获取布隆过滤器管理实例
	bloomManager := svcCtx.BloomFilters

	t.Run("UserExists", func(t *testing.T) {
		tests := []struct {
			name     string
			username string
			want     bool
			wantErr  error
		}{
			{
				name:     "Empty Username",
				username: "",
				want:     false,
				wantErr:  errorx.New(errorx.ClientError, errorx.ErrUserNotFound, errorx.Message(errorx.ErrUserNotFound)),
			},
			{
				name:     "Nonexistent Username",
				username: "nonexistent_user",
				want:     false,
				wantErr:  nil,
			},
			{
				name:     "Existing Username",
				username: "existing_user",
				want:     true,
				wantErr:  nil,
			},
		}

		ctx := context.Background()

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// 如果测试用例是 "Existing Username"，预先将用户名添加到布隆过滤器
				if tt.name == "Existing Username" {
					err := bloomManager.AddUser(ctx, tt.username)
					require.NoError(t, err, "Failed to add user to BloomFilterManager")
				}

				// 调用 UserExists 方法
				got, err := bloomManager.UserExists(ctx, tt.username)
				t.Logf("UserExists(%q) = %v, err = %v", tt.username, got, err)

				// 错误断言
				if err != nil {
					if tt.wantErr == nil {
						t.Errorf("UserExists() error = %v, wantErr %v", err, tt.wantErr)
					} else {
						// 检查错误类型和错误码
						var appErr *errorx.AppError
						ok := errors.As(err, &appErr)
						if !ok {
							t.Errorf("Expected AppError, got %T", err)
							return
						}
						var wantAppErr *errorx.AppError
						errors.As(tt.wantErr, &wantAppErr)
						if appErr.Type != wantAppErr.Type || appErr.Code != wantAppErr.Code {
							t.Errorf("UserExists() error = %v, wantErr %v", err, tt.wantErr)
						}
					}
				} else if tt.wantErr != nil {
					t.Errorf("UserExists() expected error = %v, but got none", tt.wantErr)
				}

				// 返回值断言
				if got != tt.want {
					t.Errorf("UserExists() got = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("AddUser", func(t *testing.T) {
		tests := []struct {
			name     string
			username string
			wantErr  error
		}{
			{
				name:     "Empty Username",
				username: "",
				wantErr:  errorx.New(errorx.ClientError, errorx.ErrUserNotFound, errorx.Message(errorx.ErrUserNotFound)),
			},
			{
				name:     "Valid Username",
				username: "test_user",
				wantErr:  nil,
			},
		}

		ctx := context.Background()
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := bloomManager.AddUser(ctx, tt.username)
				t.Logf("AddUser(%q) err = %v", tt.username, err)

				// 错误断言
				if err != nil {
					if tt.wantErr == nil {
						t.Errorf("AddUser() error = %v, wantErr %v", err, tt.wantErr)
					} else {
						// 检查错误类型和错误码
						var appErr *errorx.AppError
						ok := errors.As(err, &appErr)
						if !ok {
							t.Errorf("Expected AppError, got %T", err)
							return
						}
						var wantAppErr *errorx.AppError
						errors.As(tt.wantErr, &wantAppErr)
						if appErr.Type != wantAppErr.Type || appErr.Code != wantAppErr.Code {
							t.Errorf("AddUser() error = %v, wantErr %v", err, tt.wantErr)
						}
					}
				} else if tt.wantErr != nil {
					t.Errorf("AddUser() expected error = %v, but got none", tt.wantErr)
				}
			})
		}
	})
}
