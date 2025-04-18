package svc

import (
	"fmt"
	"hash/crc32"
	"strconv"
	"time"

	"gorm.io/sharding"

	"shorterurl/user/rpc/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// internal/svc/db.go
func NewDB(c config.Config, idGen func() int64) (*gorm.DB, *sharding.Sharding, error) {

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
		return nil, nil, fmt.Errorf("open database failed: %v", err)
	}

	// 3. 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, fmt.Errorf("get sql.DB failed: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 4. 配置分片中间件
	shardingInstance := sharding.Register(sharding.Config{
		ShardingKey:    c.DB.Sharding.ShardingKey,
		NumberOfShards: uint(c.DB.Sharding.NumberOfShards),
		ShardingAlgorithm: func(value any) (suffix string, err error) {
			var shard int
			switch v := value.(type) {
			case string:
				// 使用 CRC32 哈希值进行分片
				hash := crc32.ChecksumIEEE([]byte(v))
				shard = int(hash % 16) // 16 个分片
			case int, int32, int64:
				// 对整数类型的分片键取模
				var id int64
				switch val := v.(type) {
				case int:
					id = int64(val)
				case int32:
					id = int64(val)
				case int64:
					id = val
				}
				shard = int(id % 16)
			default:
				return "", fmt.Errorf("unsupported type for sharding key")
			}
			return "_" + strconv.Itoa(shard), nil // 生成 "_0" 到 "_15" 的字符串

		},
		ShardingSuffixs: func() (suffixs []string) {
			for i := 0; i < 16; i++ {
				suffixs = append(suffixs, strconv.Itoa(i))
			}
			return
		},
		PrimaryKeyGenerator: sharding.PKCustom,
		PrimaryKeyGeneratorFn: func(tableIdx int64) int64 {
			return idGen()
		},
	}, "t_user", "t_group") // 在同一个中间件中注册多个表

	// 5. 注册中间件
	if err = db.Use(shardingInstance); err != nil {
		return nil, nil, fmt.Errorf("register sharding middleware failed: %v", err)
	}

	return db, shardingInstance, nil
}
