package svc

import (
	"context"
	"flag"
	"shorterurl/admin/internal/config"
	"testing"

	"github.com/zeromicro/go-zero/core/conf"

	"github.com/stretchr/testify/require"
)

func TestBloomFilterManager_UserExists(t *testing.T) {
	configFile := flag.String("f", "../../etc/admin-api.yaml", "配置文件路径") // 定义配置文件路径标志
	flag.Parse()                                                         // 解析命令行标志

	var c config.Config            // 创建配置结构体
	conf.MustLoad(*configFile, &c) // 加载配置文件
	// 1. 初始化测试环境
	svcCtx, err := NewServiceContext(c)
	if err != nil {
		t.Fatalf("初始化服务上下文失败: %v", err)
	}

	// 2. 获取布隆过滤器实例
	bloomManager := svcCtx.BloomFilters

	//3.测试 UserExists 方法
	t.Run("UserExists", func(t *testing.T) {
		tests := []struct {
			name     string
			username string
			want     bool
			wantErr  bool
		}{
			{
				name:     "Empty Username",
				username: "",
				want:     false,
			},
			{
				name:     "Nonexistent Username",
				username: "nonexistent_user",
				want:     false,
				wantErr:  false,
			},
			{
				name:     "Existing Username",
				username: "existing_user",
				want:     true,
				wantErr:  false,
			},
		}
		ctx := context.Background()

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// 预处理：如果测试用例是 "Existing Username"，则先将用户名添加到布隆过滤器
				if tt.name == "Existing Username" {
					err := bloomManager.AddUser(ctx, tt.username)
					require.NoError(t, err, "Failed to add user to BloomFilterManager")
				}

				// 调用 UserExists 方法
				got, err := bloomManager.UserExists(ctx, tt.username)
				if (err != nil) != tt.wantErr {
					t.Logf("user := %v", got)
					t.Errorf("UserExists() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("UserExists() got = %v, want %v", got, tt.want)
				}
			})
		}
	})

	//4.测试 AddUser 方法
	t.Run("AddUser", func(t *testing.T) {
		tests := []struct {
			name     string
			username string
			wantErr  bool
		}{
			{
				name:     "Empty Username",
				username: "",
				wantErr:  true,
			},
		}
		ctx := context.Background()
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := bloomManager.AddUser(ctx, tt.username)
				if (err != nil) != tt.wantErr {
					t.Errorf("AddUser() error = %v, wantErr %v", err, tt.wantErr)
				}
			})
		}
	})
}
