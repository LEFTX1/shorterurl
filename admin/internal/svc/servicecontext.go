package svc

import (
	"go-zero-shorterurl/admin/internal/config"
	"go-zero-shorterurl/admin/internal/dal"
)

type ServiceContext struct {
	Config    config.Config
	UserRepo  *dal.UserRepo
	GroupRepo *dal.GroupRepo
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化数据库和分片配置
	if err := dal.Init(c.MySQL.DataSource); err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:    c,
		UserRepo:  dal.NewUserRepo(),
		GroupRepo: dal.NewGroupRepo(),
	}
}
