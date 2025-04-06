package svc

import (
	"shorterurl/link/rpc/shortlinkservice"
	"shorterurl/user/api/internal/config"
	"shorterurl/user/api/internal/middleware"
	"shorterurl/user/rpc/userservice"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config                  config.Config
	UserRpc                 userservice.UserService
	LinkRpc                 shortlinkservice.ShortLinkService
	Redis                   *redis.Redis
	TokenValidateMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 确保登录前缀有默认值
	if c.Auth.LoginKeyPrefix == "" {
		c.Auth.LoginKeyPrefix = "user:login:"
	}

	redisClient := redis.MustNewRedis(c.Redis.RedisConf)

	return &ServiceContext{
		Config:                  c,
		UserRpc:                 userservice.NewUserService(zrpc.MustNewClient(c.UserRpc)),
		LinkRpc:                 shortlinkservice.NewShortLinkService(zrpc.MustNewClient(c.LinkRpc)),
		Redis:                   redisClient,
		TokenValidateMiddleware: middleware.NewTokenValidateMiddleware(&c.Auth, redisClient).Handle,
	}
}
