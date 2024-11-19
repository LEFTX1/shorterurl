package svc

import (
	"fmt"
	"gorm.io/sharding"

	"go-zero-shorterurl/admin/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 修改NewDB函数签名，接收ID生成器
func NewDB(c config.Config, idGen func() int64) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DB.User, c.DB.Password, c.DB.Host, c.DB.Port, c.DB.Database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open database failed: %v", err)
	}

	// 配置分片中间件
	middleware := sharding.Register(sharding.Config{
		ShardingKey:         "username",
		NumberOfShards:      16,
		PrimaryKeyGenerator: sharding.PKCustom,
		PrimaryKeyGeneratorFn: func(tableIdx int64) int64 {
			return idGen()
		},
	}, "t_user")

	if err = db.Use(middleware); err != nil {
		return nil, fmt.Errorf("register sharding middleware failed: %v", err)
	}

	return db, nil
}
