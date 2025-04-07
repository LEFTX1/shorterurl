package consumer

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
)

// ServiceContext 服务上下文接口
type ServiceContext interface {
	GetRedis() *redis.Redis
	GetDBs() DBInterface
}

// DBInterface 数据库接口
type DBInterface interface {
	GetCommon() *gorm.DB
	GetLinkDB() *gorm.DB
}

// StatsConsumer 统计消费者接口
type StatsConsumer interface {
	Start()
	Stop()
	Submit(record *StatsRecord)
}
