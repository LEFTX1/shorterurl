package svc

import (
	"fmt"
	"go-zero-shorterurl/pkg/snowflake"
)

var idGenerator func() int64

// InitSnowflake 初始化服务的雪花算法生成器
func InitSnowflake() error {
	// 初始化雪花算法节点池
	if err := snowflake.InitSnowflake(); err != nil {
		return fmt.Errorf("init snowflake failed: %v", err)
	}

	// 获取ID生成器
	generator, err := snowflake.GetSnowflakeGenerator()
	if err != nil {
		return fmt.Errorf("get snowflake generator failed: %v", err)
	}

	idGenerator = generator
	return nil
}

// GetIDGenerator 获取ID生成器
func GetIDGenerator() func() int64 {
	return idGenerator
}
