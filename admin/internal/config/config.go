package config

import (
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	DB struct {
		Host     string
		Port     int
		User     string
		Password string
		Database string
	}
}
