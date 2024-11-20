package svc

import (
	"fmt"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"go-zero-shorterurl/admin/internal/config"
	"go-zero-shorterurl/admin/internal/dal/query"
	"go-zero-shorterurl/pkg/snowflake"
	"gorm.io/gorm"
)

// internal/svc/servicecontext.go
type ServiceContext struct {
	Config       config.Config
	DB           *gorm.DB
	Query        *query.Query
	Redis        *RedisClient
	BloomFilters *BloomFilterManager
	RS           *redsync.Redsync
}

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

	// 3. 初始化MySQL
	db, err := NewDB(c, idGen)
	if err != nil {
		return nil, fmt.Errorf("init database failed: %v", err)
	}

	// 4. 初始化查询对象
	q := query.Use(db)

	// 5. 初始化Redis客户端
	redis, err := NewRedisClient(c.Redis.Host, c.Redis.Password, c.Redis.DB)
	if err != nil {
		return nil, fmt.Errorf("init redis failed: %v", err)
	}

	// 6. 初始化布隆过滤器管理器
	bloomFilters, err := NewBloomFilterManager(redis.Client, c)
	if err != nil {
		return nil, fmt.Errorf("init bloom filters failed: %v", err)
	}

	// 7. 初始化分布式锁
	pool := goredis.NewPool(redis.Client)
	rs := redsync.New(pool)

	return &ServiceContext{
		Config:       c,
		DB:           db,
		Query:        q,
		Redis:        redis,
		BloomFilters: bloomFilters,
		RS:           rs,
	}, nil
}

// 添加 Close 方法
func (s *ServiceContext) Close() error {
	// 1. 关闭数据库连接
	sqlDB, err := s.DB.DB()
	if err != nil {
		return fmt.Errorf("get sql.DB failed: %v", err)
	}
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("close database connection failed: %v", err)
	}

	// 2. 关闭Redis连接
	if err := s.Redis.Client.Close(); err != nil {
		return fmt.Errorf("close redis connection failed: %v", err)
	}

	return nil
}
