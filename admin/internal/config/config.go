package config

import (
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
		Host        string
		Password    string
		DB          int
		BloomFilter struct {
			User struct {
				Name            string  `json:",default=user:register:filter"`
				ErrorRate       float64 `json:",default=0.01"`
				InitialCapacity int64   `json:",default=1000000"`
			}
			Group struct {
				Name            string  `json:",default=group:register:filter"`
				ErrorRate       float64 `json:",default=0.01"`
				InitialCapacity int64   `json:",default=2000000"`
			}
		}
	}
}
