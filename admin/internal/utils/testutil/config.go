// internal/utils/testutil/config.go
package testutil

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"go-zero-shorterurl/admin/internal/config"
	"go-zero-shorterurl/admin/internal/svc"
	"path/filepath"
	"runtime"
)

// GetTestConfig 获取测试配置和服务上下文
func GetTestConfig() (*svc.ServiceContext, error) {
	// 1. 获取测试配置文件路径
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("get caller info failed")
	}

	configFile := filepath.Join(filepath.Dir(filename), "../../../etc/admin-api.admin-api.test.yaml")

	// 2. 加载配置
	var c config.Config
	conf.MustLoad(configFile, &c)

	// 3. 创建服务上下文
	return svc.NewServiceContext(c)
}
