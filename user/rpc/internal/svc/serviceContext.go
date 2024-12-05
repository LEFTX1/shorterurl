package svc

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
	"gorm.io/sharding"
	"shorterurl/user/rpc/internal/common"
	"shorterurl/user/rpc/internal/config"
	"shorterurl/user/rpc/internal/dal/query"
	"shorterurl/user/rpc/pkg/snowflake"
)

type ServiceContext struct {
	Config       config.Config
	DB           *gorm.DB
	Query        *query.Query
	Redis        *redis.Redis
	BloomFilters *BloomFilterManager
	Sharding     *sharding.Sharding
}

func NewServiceContext(c config.Config) *ServiceContext {
	// Initialize snowflake
	if err := snowflake.InitSnowflake(); err != nil {
		panic(fmt.Errorf("init snowflake failed: %v", err))
		return nil
	}

	// Initialize AES encryption
	if err := common.InitAES(c); err != nil {
		panic(fmt.Errorf("init AES failed: %v", err))
		return nil
	}

	// Get ID generator
	idGen, err := snowflake.GetSnowflakeGenerator()
	if err != nil {
		panic(fmt.Errorf("get snowflake generator failed: %v", err))
		return nil
	}

	// Initialize MySQL and get sharding instance
	db, shardingInstance, err := NewDB(c, idGen)
	if err != nil {
		panic(fmt.Errorf("init database failed: %v", err))
		return nil
	}

	// Initialize query object
	q := query.Use(db)

	// Initialize Redis client
	redisClient := redis.MustNewRedis(c.BizRedis)

	// Initialize Bloom filter manager
	bloomFilters := NewBloomFilterManager(redisClient, c)

	return &ServiceContext{
		Config:       c,
		DB:           db,
		Query:        q,
		Redis:        redisClient,
		BloomFilters: bloomFilters,
		Sharding:     shardingInstance,
	}
}
