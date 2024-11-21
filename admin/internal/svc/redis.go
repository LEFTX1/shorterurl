// internal/svc/redis.go
package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go-zero-shorterurl/admin/internal/config"
)

type RedisClient struct {
	*redis.Redis
}

func NewRedisClient(c config.Config) (*RedisClient, error) {
	// 使用 go-zero 的 Redis 客户端
	client := redis.MustNewRedis(c.Redis.RedisConf)
	return &RedisClient{Redis: client}, nil
}
