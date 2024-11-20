// internal/svc/redis.go
package svc

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

// RedisClient Redis客户端封装
type RedisClient struct {
	Client *redis.Client // 原生redis客户端实例
}

// NewRedisClient 创建新的Redis客户端
// addr: Redis服务器地址，格式为 host:port
// password: Redis密码，如果没有则传空字符串
// db: Redis数据库编号
func NewRedisClient(addr, password string, db int) (*RedisClient, error) {
	// 创建Redis客户端配置
	client := redis.NewClient(&redis.Options{
		Addr:         addr,            // Redis地址
		Password:     password,        // Redis密码
		DB:           db,              // 数据库编号
		DialTimeout:  5 * time.Second, // 连接超时
		ReadTimeout:  3 * time.Second, // 读取超时
		WriteTimeout: 3 * time.Second, // 写入超时
		PoolSize:     10,              // 连接池大小
		PoolTimeout:  4 * time.Second, // 连接池超时
	})

	// 测试连接是否成功
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err // 连接失败，返回错误
	}

	return &RedisClient{Client: client}, nil
}
