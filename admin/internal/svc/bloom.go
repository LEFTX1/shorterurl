// internal/svc/bloom.go
package svc

import (
	"context"
	"shorterurl/admin/internal/config"
	"shorterurl/admin/internal/types/errorx"

	"github.com/zeromicro/go-zero/core/bloom"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type BloomFilterManager struct {
	userFilter  *bloom.Filter
	groupFilter *bloom.Filter
}

func NewBloomFilterManager(redisClient *redis.Redis, c config.Config) *BloomFilterManager {
	return &BloomFilterManager{
		userFilter:  bloom.New(redisClient, c.Redis.BloomFilter.User.Name, c.Redis.BloomFilter.User.Bits),
		groupFilter: bloom.New(redisClient, c.Redis.BloomFilter.Group.Name, c.Redis.BloomFilter.Group.Bits),
	}
}

// 用户相关方法
func (m *BloomFilterManager) AddUser(ctx context.Context, username string) error {
	if username == "" {
		return errorx.NewUserError(errorx.UserNotExistError)
	}
	return m.userFilter.AddCtx(ctx, []byte(username))
}

func (m *BloomFilterManager) UserExists(ctx context.Context, username string) (bool, error) {
	if username == "" {
		return false, errorx.NewUserError(errorx.UserNotExistError)
	}
	return m.userFilter.ExistsCtx(ctx, []byte(username))
}

// 群组相关方法
func (m *BloomFilterManager) AddGroup(ctx context.Context, groupname string) error {
	return m.groupFilter.AddCtx(ctx, []byte(groupname))
}

func (m *BloomFilterManager) GroupExists(ctx context.Context, groupname string) (bool, error) {
	return m.groupFilter.ExistsCtx(ctx, []byte(groupname))
}
