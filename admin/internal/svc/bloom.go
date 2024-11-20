// internal/svc/bloom.go
package svc

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go-zero-shorterurl/admin/internal/config"
)

// BloomFilter 布隆过滤器
type BloomFilter struct {
	rdb *redis.Client
	key string
}

// BloomFilterManager 布隆过滤器管理器
type BloomFilterManager struct {
	UserFilter  *BloomFilter // 用户名布隆过滤器
	GroupFilter *BloomFilter // 分组名布隆过滤器
}

// NewBloomFilter 创建单个布隆过滤器
func NewBloomFilter(rdb *redis.Client, key string) *BloomFilter {
	return &BloomFilter{
		rdb: rdb,
		key: key,
	}
}

// Init 初始化布隆过滤器
func (bf *BloomFilter) Init(ctx context.Context, capacity int64, errorRate float64) error {
	// 使用 BF.RESERVE 命令初始化布隆过滤器
	// 如果已存在会返回错误："ERR item exists"，这种情况可以忽略
	err := bf.rdb.Do(ctx, "BF.RESERVE", bf.key, errorRate, capacity).Err()
	if err != nil && err.Error() != "ERR item exists" {
		return fmt.Errorf("init bloom filter failed: %v", err)
	}
	return nil
}

// Add 添加元素到布隆过滤器
func (bf *BloomFilter) Add(ctx context.Context, item string) error {
	return bf.rdb.Do(ctx, "BF.ADD", bf.key, item).Err()
}

// Exists 检查元素是否存在
func (bf *BloomFilter) Exists(ctx context.Context, item string) (bool, error) {
	return bf.rdb.Do(ctx, "BF.EXISTS", bf.key, item).Bool()
}

// NewBloomFilterManager 创建布隆过滤器管理器
func NewBloomFilterManager(rdb *redis.Client, config config.Config) (*BloomFilterManager, error) {
	ctx := context.Background()

	// 初始化用户名布隆过滤器
	userFilter := NewBloomFilter(rdb, config.Redis.BloomFilter.User.Name)
	if err := userFilter.Init(ctx,
		config.Redis.BloomFilter.User.InitialCapacity,
		config.Redis.BloomFilter.User.ErrorRate,
	); err != nil {
		return nil, fmt.Errorf("init user bloom filter failed: %v", err)
	}

	// 初始化分组名布隆过滤器
	groupFilter := NewBloomFilter(rdb, config.Redis.BloomFilter.Group.Name)
	if err := groupFilter.Init(ctx,
		config.Redis.BloomFilter.Group.InitialCapacity,
		config.Redis.BloomFilter.Group.ErrorRate,
	); err != nil {
		return nil, fmt.Errorf("init group bloom filter failed: %v", err)
	}

	return &BloomFilterManager{
		UserFilter:  userFilter,
		GroupFilter: groupFilter,
	}, nil
}

// AddUser 添加用户名到布隆过滤器
func (m *BloomFilterManager) AddUser(ctx context.Context, username string) error {
	return m.UserFilter.Add(ctx, username)
}

// UserExists 检查用户名是否可能存在
func (m *BloomFilterManager) UserExists(ctx context.Context, username string) (bool, error) {
	return m.UserFilter.Exists(ctx, username)
}

// AddGroup 添加分组名到布隆过滤器
func (m *BloomFilterManager) AddGroup(ctx context.Context, groupName string) error {
	return m.GroupFilter.Add(ctx, groupName)
}

// GroupExists 检查分组名是否可能存在
func (m *BloomFilterManager) GroupExists(ctx context.Context, groupName string) (bool, error) {
	return m.GroupFilter.Exists(ctx, groupName)
}
