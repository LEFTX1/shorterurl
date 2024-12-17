package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	DB struct {
		Host     string
		Port     int
		User     string
		Password string
		Database string
		Sharding struct {
			ShardingKey    string `json:",default=gid"`
			NumberOfShards int    `json:",default=16"`
		}
	}

	BizRedis redis.RedisConf

	BloomFilter struct {
		Name         string `json:",default=test:bloom:shortlinks"`
		Size         uint   `json:",default=20000000"`
		RedisKeyName string `json:",default=bloom:shortlinks"`
	}

	Domain string `json:",default=s.xleft.cn"`
}
