package svc

import (
	"fmt"
	"shorterurl/link/rpc/internal/config"
	"shorterurl/link/rpc/internal/dal/query"
	"shorterurl/link/rpc/pkg/snowflake"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
	"gorm.io/sharding"
)

type ServiceContext struct {
	Config         config.Config
	DB             *gorm.DB
	Query          *query.Query
	BizRedis       *redis.Redis
	BloomFilterMgr *BloomFilterManager
	Sharding       *sharding.Sharding
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化雪花算法
	if err := snowflake.InitSnowflake(); err != nil {
		panic(fmt.Errorf("init snowflake failed: %v", err))
	}

	// 获取ID生成器
	idGen, err := snowflake.GetSnowflakeGenerator()
	if err != nil {
		panic(fmt.Errorf("init id generator failed: %v", err))
	}

	// 初始化数据库连接
	db, shardingInstance, err := InitDB(c, idGen)
	if err != nil {
		panic(fmt.Errorf("init database failed: %v", err))
	}

	// 初始化Redis客户端
	bizRedis := redis.MustNewRedis(c.BizRedis)

	// 初始化布隆过滤器管理器
	bloomFilterMgr, err := NewBloomFilterManager(bizRedis, c)
	if err != nil {
		panic(fmt.Errorf("init bloom filter failed: %v", err))
	}

	return &ServiceContext{
		Config:         c,
		DB:             db,
		Query:          query.Use(db),
		BizRedis:       bizRedis,
		BloomFilterMgr: bloomFilterMgr,
		Sharding:       shardingInstance,
	}
}
