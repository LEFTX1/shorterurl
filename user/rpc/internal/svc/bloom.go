// internal/svc/bloom.go
package svc

import (
	"context"
	"shorterurl/user/rpc/internal/config"
	"shorterurl/user/rpc/internal/types/errorx"

	"github.com/zeromicro/go-zero/core/bloom"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type BloomFilterManager struct {
	userFilter  *bloom.Filter
	groupFilter *bloom.Filter
}

func NewBloomFilterManager(redisClient *redis.Redis, c config.Config) *BloomFilterManager {
	return &BloomFilterManager{
		userFilter:  bloom.New(redisClient, c.BloomFilter.User.Name, c.BloomFilter.User.Bits),
		groupFilter: bloom.New(redisClient, c.BloomFilter.Group.Name, c.BloomFilter.Group.Bits),
	}
}

// AddUser 用户相关方法
func (m *BloomFilterManager) AddUser(ctx context.Context, username string) error {
	if username == "" {
		return errorx.New(errorx.ClientError, errorx.ErrUserNotFound, errorx.Message(errorx.ErrUserNotFound))
	}
	return m.userFilter.AddCtx(ctx, []byte(username))
}

func (m *BloomFilterManager) UserExists(ctx context.Context, username string) (bool, error) {
	if username == "" {
		return false, errorx.New(errorx.ClientError, errorx.ErrUserNotFound, errorx.Message(errorx.ErrUserNotFound))
	}
	return m.userFilter.ExistsCtx(ctx, []byte(username))
}

// AddGroup 群组相关方法
func (m *BloomFilterManager) AddGroup(ctx context.Context, groupname string) error {
	if groupname == "" {
		return errorx.New(errorx.ClientError, errorx.ErrGroupNotFound, errorx.Message(errorx.ErrGroupNotFound))
	}
	return m.groupFilter.AddCtx(ctx, []byte(groupname))
}

func (m *BloomFilterManager) GroupExists(ctx context.Context, groupname string) (bool, error) {
	if groupname == "" {
		return false, errorx.New(errorx.ClientError, errorx.ErrGroupNotFound, errorx.Message(errorx.ErrGroupNotFound))
	}
	return m.groupFilter.ExistsCtx(ctx, []byte(groupname))
}
