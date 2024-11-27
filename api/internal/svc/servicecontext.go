package svc

import (
	"shorterurl/api/internal/config"
	"shorterurl/api/internal/middleware"
	"shorterurl/api/internal/user_client"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config         config.Config
	Redis          *redis.Redis
	TokenValidator *middleware.TokenValidateMiddleware
	UserRpc        user_client.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	redisClient := redis.MustNewRedis(c.Redis.RedisConf)
	return &ServiceContext{
		Config:         c,
		Redis:          redisClient,
		TokenValidator: middleware.NewTokenValidateMiddleware(&c.Auth, redisClient),
		UserRpc:        user_client.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}
