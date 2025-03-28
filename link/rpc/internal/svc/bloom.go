package svc

import (
	"context"
	"shorterurl/link/rpc/internal/config"
	"shorterurl/link/rpc/internal/types/errorx"

	"github.com/zeromicro/go-zero/core/bloom"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// BloomFilterManager 布隆过滤器管理器
type BloomFilterManager struct {
	filter *bloom.Filter
	redis  *redis.Redis
	config config.Config
}

// NewBloomFilterManager 创建新的布隆过滤器管理器
func NewBloomFilterManager(redisClient *redis.Redis, c config.Config) (*BloomFilterManager, error) {
	if redisClient == nil {
		return nil, errorx.NewCodeError(errorx.ErrRedisClientNil, errorx.ErrBloomFilterInit, errorx.MsgRedisNotFound)
	}

	// 创建布隆过滤器
	filter := bloom.New(redisClient, c.BloomFilter.RedisKeyName, c.BloomFilter.Size)
	if filter == nil {
		return nil, errorx.NewCodeError(errorx.ErrBloomFilterInit, errorx.ErrBloomFilterInit, errorx.MsgBloomFilterInit)
	}

	return &BloomFilterManager{
		filter: filter,
		redis:  redisClient,
		config: c,
	}, nil
}

// Add 添加短链接到布隆过滤器
func (m *BloomFilterManager) Add(ctx context.Context, shortLink string) error {
	if shortLink == "" {
		return errorx.NewCodeError(errorx.ErrShortLinkEmpty, errorx.ErrBloomFilterAdd, errorx.MsgShortLinkEmpty)
	}

	if err := m.filter.Add([]byte(shortLink)); err != nil {
		return errorx.NewCodeError(errorx.ErrBloomFilterAdd, errorx.ErrBloomFilterAdd,
			"%s: %v", errorx.MsgBloomFilterAdd, err)
	}

	return nil
}

// Exists 检查短链接是否存在于布隆过滤器中
func (m *BloomFilterManager) Exists(ctx context.Context, shortLink string) (bool, error) {
	if shortLink == "" {
		return false, errorx.NewCodeError(errorx.ErrShortLinkEmpty, errorx.ErrBloomFilterCheck, errorx.MsgShortLinkEmpty)
	}

	exists, err := m.filter.Exists([]byte(shortLink))
	if err != nil {
		return false, errorx.NewCodeError(errorx.ErrBloomFilterCheck, errorx.ErrBloomFilterCheck,
			"%s: %v", errorx.MsgBloomFilterCheck, err)
	}

	return exists, nil
}

// Reset 重置布隆过滤器
func (m *BloomFilterManager) Reset() error {
	// 由于bloom.Filter没有直接提供Reset方法
	// 我们通过删除原有的key并重新创建过滤器来实现
	key := m.config.BloomFilter.RedisKeyName

	// 删除Redis中的布隆过滤器键
	_, err := m.redis.Del(key)
	if err != nil {
		return errorx.NewCodeError(errorx.ErrBloomFilterReset, errorx.ErrBloomFilterReset,
			"%s: %v", errorx.MsgBloomFilterReset, err)
	}

	// 重新创建布隆过滤器
	newFilter := bloom.New(m.redis, key, m.config.BloomFilter.Size)
	if newFilter == nil {
		return errorx.NewCodeError(errorx.ErrBloomFilterInit, errorx.ErrBloomFilterReset, errorx.MsgBloomFilterInit)
	}

	// 更新过滤器
	m.filter = newFilter

	return nil
}

// BatchAdd 批量添加短链接到布隆过滤器
func (m *BloomFilterManager) BatchAdd(ctx context.Context, shortLinks []string) error {
	if len(shortLinks) == 0 {
		return nil
	}

	for _, link := range shortLinks {
		if err := m.Add(ctx, link); err != nil {
			return err
		}
	}

	return nil
}

// BatchExists 批量检查短链接是否存在
func (m *BloomFilterManager) BatchExists(ctx context.Context, shortLinks []string) (map[string]bool, error) {
	result := make(map[string]bool)

	if len(shortLinks) == 0 {
		return result, nil
	}

	for _, link := range shortLinks {
		exists, err := m.Exists(ctx, link)
		if err != nil {
			return nil, err
		}
		result[link] = exists
	}

	return result, nil
}
