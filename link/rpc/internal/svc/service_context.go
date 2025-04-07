package svc

import (
	"fmt"
	"shorterurl/link/rpc/internal/config"
	"shorterurl/link/rpc/internal/consumer"
	"shorterurl/link/rpc/internal/repo"
	"shorterurl/link/rpc/pkg/snowflake"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

// ServiceContext 服务上下文
type ServiceContext struct {
	Config         config.Config
	DBs            *DBs
	BizRedis       *redis.Redis
	BloomFilterMgr *BloomFilterManager
	RepoManager    *repo.RepoManager
	StatsConsumer  *consumer.ShortLinkStatsConsumer
}

// 实现消费者所需的接口
func (s *ServiceContext) GetRedis() *redis.Redis {
	return s.BizRedis
}

func (s *ServiceContext) GetDBs() consumer.DBInterface {
	return s.DBs
}

// NewServiceContext 创建服务上下文
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
	dbs, err := InitDBs(c, idGen)
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

	// 初始化仓库管理器
	repoManager := repo.NewRepoManager(
		dbs.Common,
		dbs.LinkDB,
		dbs.GotoLinkDB,
		dbs.GroupDB,
		dbs.UserDB,
	)

	svcCtx := &ServiceContext{
		Config:         c,
		DBs:            dbs,
		BizRedis:       bizRedis,
		BloomFilterMgr: bloomFilterMgr,
		RepoManager:    repoManager,
	}

	// 创建并启动统计消费者
	statsConsumer := consumer.NewShortLinkStatsConsumer(svcCtx)
	statsConsumer.Start()
	svcCtx.StatsConsumer = statsConsumer

	return svcCtx
}
