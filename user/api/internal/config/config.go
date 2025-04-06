package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Auth TokenValidateConfig `json:"Auth"`

	Redis struct {
		RedisConf redis.RedisConf
	}

	UserRpc zrpc.RpcClientConf
	LinkRpc zrpc.RpcClientConf
}
