package svc

import (
	"fmt"
	"go-zero-shorterurl/admin/internal/config"
	"go-zero-shorterurl/admin/internal/dal/query"
	"go-zero-shorterurl/pkg/snowflake"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
	Query  *query.Query
}

func NewServiceContext(c config.Config) (*ServiceContext, error) {
	// 1. 初始化雪花算法
	if err := snowflake.InitSnowflake(); err != nil {
		return nil, fmt.Errorf("init snowflake failed: %v", err)
	}

	// 2. 获取ID生成器（仅用于初始化DB）
	idGen, err := snowflake.GetSnowflakeGenerator()
	if err != nil {
		return nil, fmt.Errorf("get snowflake generator failed: %v", err)
	}

	// 3. 初始化数据库连接
	db, err := NewDB(c, idGen)
	if err != nil {
		return nil, fmt.Errorf("init database failed: %v", err)
	}

	// 4. 初始化查询对象
	q := query.Use(db)

	// 5. 返回服务上下文（不包含ID生成器）
	return &ServiceContext{
		Config: c,
		DB:     db,
		Query:  q,
	}, nil
}
