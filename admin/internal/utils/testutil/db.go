package testutil

import (
	"fmt"
	"log"
	"shorterurl/admin/pkg/snowflake"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/sharding"
)

var (
	once    sync.Once
	idGen   func() int64
	initErr error
)

// TestDBConfig 测试数据库配置
type TestDBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// GetIDGenerator 获取ID生成器
func GetIDGenerator() (func() int64, error) {
	once.Do(func() {
		// 初始化雪花算法
		if err := snowflake.InitSnowflake(); err != nil {
			initErr = fmt.Errorf("init snowflake failed: %v", err)
			return
		}

		// 获取ID生成器
		generator, err := snowflake.GetSnowflakeGenerator()
		if err != nil {
			initErr = fmt.Errorf("get snowflake generator failed: %v", err)
			return
		}

		idGen = generator
	})

	if initErr != nil {
		return nil, initErr
	}

	return idGen, nil
}

// GetDefaultTestConfig 获取默认测试配置
func GetDefaultTestConfig() TestDBConfig {
	return TestDBConfig{
		Host:     "localhost",
		Port:     3306,
		User:     "root",
		Password: "123456",
		Database: "link_go",
	}
}

// GetTestDB 获取测试用的DB实例
func GetTestDB() (*gorm.DB, error) {
	// 获取ID生成器
	generator, err := GetIDGenerator()
	if err != nil {
		return nil, fmt.Errorf("get id generator failed: %v", err)
	}

	// 使用默认配置
	config := GetDefaultTestConfig()

	// 构建DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	// 连接数据库
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
			return generator()
		},
	}, "t_user") // 确认这里是正确的表名

	// 打印分片配置
	log.Printf("分片配置: 表=%s, 分片数=%d", "t_user", 16)

	// 注册分片中间件
	if err = db.Use(middleware); err != nil {
		return nil, fmt.Errorf("register sharding middleware failed: %v", err)
	}

	return db, nil
}
