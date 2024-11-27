// internal/svc/bloom.go
package svc

import (
	"shorterurl/rpc/user/internal/config"

	"github.com/zeromicro/go-zero/core/bloom"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// BloomFilterManager 布隆过滤器管理器
type BloomFilterManager struct {
	userFilter *bloom.Filter
}

// NewBloomFilterManager 创建布隆过滤器管理器
func NewBloomFilterManager(store *redis.Redis, c config.Config) *BloomFilterManager {
	return &BloomFilterManager{
		userFilter: bloom.New(store, "user_filter", 20000000), // 预计用户量为 2000w
	}
}

// UserExists 检查用户是否存在
func (m *BloomFilterManager) UserExists(username string) bool {
	return m.userFilter.Exists([]byte(username))
}

// AddUser 添加用户到布隆过滤器
func (m *BloomFilterManager) AddUser(username string) error {
	return m.userFilter.Add([]byte(username))
}
