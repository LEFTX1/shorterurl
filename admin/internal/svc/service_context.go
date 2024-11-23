package svc

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go-zero-shorterurl/admin/internal/config"
	"go-zero-shorterurl/admin/internal/dal/query"
	"go-zero-shorterurl/pkg/snowflake"
	"gorm.io/gorm"
	"gorm.io/sharding"
)

// internal/svc/service_context.go
type ServiceContext struct {
	Config       config.Config
	DB           *gorm.DB
	Query        *query.Query
	Redis        *redis.Redis // 使用 go-zero 的 Redis 客户端
	BloomFilters *BloomFilterManager
	Sharding     *sharding.Sharding // 添加分片实例
}

// internal/svc/service_context.go
func NewServiceContext(c config.Config) (*ServiceContext, error) {
	// 1. 初始化雪花算法
	if err := snowflake.InitSnowflake(); err != nil {
		return nil, fmt.Errorf("init snowflake failed: %v", err)
	}

	// 2. 获取ID生成器
	idGen, err := snowflake.GetSnowflakeGenerator()
	if err != nil {
		return nil, fmt.Errorf("get snowflake generator failed: %v", err)
	}

	// 3. 初始化MySQL并获取分片实例
	db, shardingInstance, err := NewDB(c, idGen)
	if err != nil {
		return nil, fmt.Errorf("init database failed: %v", err)
	}

	// 4. 初始化查询对象
	q := query.Use(db)

	// 5. 初始化Redis客户端
	redisClient := redis.MustNewRedis(c.Redis.RedisConf)

	// 6. 初始化布隆过滤器管理器
	bloomFilters := NewBloomFilterManager(redisClient, c)

	return &ServiceContext{
		Config:       c,
		DB:           db,
		Query:        q,
		Redis:        redisClient,
		BloomFilters: bloomFilters,
		Sharding:     shardingInstance, // 添加分片实例
	}, nil
}
