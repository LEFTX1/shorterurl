package svc

import (
	"fmt"
	"gorm.io/sharding"
	"time"

	"go-zero-shorterurl/admin/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// internal/svc/db.go
func NewDB(c config.Config, idGen func() int64) (*gorm.DB, error) {
	// 1. 构建DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DB.User,
		c.DB.Password,
		c.DB.Host,
		c.DB.Port,
		c.DB.Database,
	)

	// 2. 打开数据库连接
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open database failed: %v", err)
	}

	// 3. 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql.DB failed: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 4. 配置分片中间件
	middleware := sharding.Register(sharding.Config{
		ShardingKey:         c.DB.Sharding.ShardingKey,
		NumberOfShards:      uint(c.DB.Sharding.NumberOfShards),
		PrimaryKeyGenerator: sharding.PKCustom,
		PrimaryKeyGeneratorFn: func(tableIdx int64) int64 {
			return idGen()
		},
	}, "t_user", "t_group") // 在同一个中间件中注册多个表

	// 5. 注册中间件
	if err = db.Use(middleware); err != nil {
		return nil, fmt.Errorf("register sharding middleware failed: %v", err)
	}

	return db, nil
}
