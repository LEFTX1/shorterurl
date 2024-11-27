package svc

import (
	"fmt"
	"shorterurl/rpc/user/internal/config"
	"shorterurl/rpc/user/pkg/snowflake"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/sharding"
)

// NewDB 初始化数据库连接和分片配置
func NewDB(c config.Config, idGen *snowflake.Generator) (*gorm.DB, *sharding.Sharding, error) {
	// 连接数据库
	db, err := gorm.Open(mysql.Open(c.DB.DataSource), &gorm.Config{})
	if err != nil {
		return nil, nil, fmt.Errorf("connect database failed: %v", err)
	}

	// 配置分片规则
	shardingInstance := sharding.Register(sharding.Config{
		ShardingKey:    "username",
		NumberOfShards: 16,
		PrimaryKeyGenerator: func(shardingValue interface{}) (int64, error) {
			return idGen.NextID()
		},
	}, db)

	return db, shardingInstance, nil
}
