// internal/utils/testutil/config.go
package testutil

// TestDBConfig 测试数据库配置
type TestDBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// GetTestDBConfig 获取测试配置
func GetTestDBConfig() TestDBConfig {
	return TestDBConfig{
		Host:     "127.0.0.1",
		Port:     3306,
		User:     "root",
		Password: "123456",
		Database: "link_go", // 修改为 link 数据库
	}
}
