package dal

import (
	"fmt"
	"go-zero-shorterurl/admin/internal/dal/query"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/sharding"
	"sync"
)

var (
	DB          *gorm.DB
	once        sync.Once
	initialized bool
)

// Init 初始化数据库和分片配置
func Init(dsn string) error {
	var err error

	// 使用 once 确保只初始化一次
	once.Do(func() {
		if initialized {
			return
		}

		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			err = fmt.Errorf("connect db error: %w", err)
			return
		}

		// 用户表分片中间件
		userMiddleware := sharding.Register(sharding.Config{
			ShardingKey:         "username",
			NumberOfShards:      16,
			PrimaryKeyGenerator: sharding.PKSnowflake,
		}, "t_user")

		// 分组表分片中间件
		groupMiddleware := sharding.Register(sharding.Config{
			ShardingKey:         "username",
			NumberOfShards:      16,
			PrimaryKeyGenerator: sharding.PKSnowflake,
		}, "t_group")

		// 使用分片中间件
		DB.Use(userMiddleware)
		DB.Use(groupMiddleware)

		initialized = true
	})

	return err
}

// GetQuery 获取查询实例
func GetQuery() *query.Query {
	return query.Use(DB)
}
