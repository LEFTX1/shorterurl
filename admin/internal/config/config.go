package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

// internal/config/config.go
type Config struct {
	rest.RestConf
	DB struct {
		Host     string
		Port     int
		User     string
		Password string
		Database string
		Sharding struct {
			ShardingKey    string `json:",default=username"`
			NumberOfShards int    `json:",default=16"`
		}
	}
	Redis struct {
		RedisConf   redis.RedisConf
		BloomFilter struct {
			User struct {
				Name            string `json:",default=bloom:users"`
				Bits            uint   `json:",default=20000000"`
				InitialCapacity int64  `json:",default=1000000"`
			}
			Group struct {
				Name            string `json:",default=bloom:groups"`
				Bits            uint   `json:",default=20000000"`
				InitialCapacity int64  `json:",default=1000000"`
			}
		}
	}
	Crypto struct {
		AESKey string `json:"aesKey"` // base64 编码的 AES 密钥
	}
}
