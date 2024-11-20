package testutil

import (
	"fmt"
	"gorm.io/sharding"

	"go-zero-shorterurl/pkg/snowflake"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"sync"
)

var (
	once    sync.Once
	idGen   func() int64
	initErr error
)

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

// GetTestDB 获取测试用的DB实例
func GetTestDB() (*gorm.DB, error) {
	// 获取ID生成器
	generator, err := GetIDGenerator()
	if err != nil {
		return nil, fmt.Errorf("get id generator failed: %v", err)
	}

	conf := GetTestDBConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.User, conf.Password, conf.Host, conf.Port, conf.Database)

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

// CleanTestData 清理测试数据
func CleanTestData(db *gorm.DB) error {
	for i := 0; i < 16; i++ {
		// 使用 %d 而不是 %02d，这样会生成 t_user_0, t_user_1 等格式
		tableName := fmt.Sprintf("t_user_%d", i)
		if err := db.Exec(fmt.Sprintf("DELETE FROM %s WHERE username LIKE 'test_%%'", tableName)).Error; err != nil {
			log.Printf("Warning: Failed to clean table %s: %v", tableName, err)
		}
	}
	return nil
}
