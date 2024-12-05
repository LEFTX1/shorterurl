package logic

import (
	"context"
	"math/rand"
	"shorterurl/user/rpc/internal/constant"
	"shorterurl/user/rpc/internal/dal/model"
	"shorterurl/user/rpc/internal/svc"
	"shorterurl/user/rpc/internal/types/errorx"
	__ "shorterurl/user/rpc/pb"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

const (
	totalCharacters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	base            = 62 // len(totalCharacters)
)

type GroupCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupCreateLogic {
	return &GroupCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// generateRandomString 生成随机字符串
func generateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = totalCharacters[rand.Intn(base)]
	}
	return string(b)
}

// GroupCreate 创建分组
func (l *GroupCreateLogic) GroupCreate(in *__.GroupSaveRequest) (*__.CommonResponse, error) {
	// 1. 创建分布式锁，防止并发创建相同分组
	lockKey := constant.LOCK_GROUP_CREATE_KEY + in.Username
	lock := redis.NewRedisLock(l.svcCtx.Redis, lockKey)
	lock.SetExpire(30) // 设置锁的过期时间为30秒

	// 2. 获取锁
	acquired, err := lock.AcquireCtx(l.ctx)
	if err != nil {
		return nil, errorx.New(errorx.SystemError, errorx.ErrDistributedLock, errorx.Message(errorx.ErrDistributedLock))
	}
	if !acquired {
		return nil, errorx.New(errorx.ClientError, errorx.ErrTooManyRequests, errorx.Message(errorx.ErrTooManyRequests))
	}

	// 确保在函数结束时释放锁
	defer func() {
		if released, err := lock.Release(); err != nil {
			logx.Errorf("释放锁失败: %v", err)
		} else if !released {
			logx.Error("锁未被主动释放")
		}
	}()

	// 3. 查询当前用户的所有分组
	groupCount, err := l.svcCtx.Query.TGroup.WithContext(l.ctx).
		Where(l.svcCtx.Query.TGroup.Username.Eq(in.Username)).
		Where(l.svcCtx.Query.TGroup.DelFlag.Is(false)).
		Count()
	if err != nil {
		return nil, errorx.New(errorx.SystemError, errorx.ErrInternalServer, errorx.Message(errorx.ErrInternalServer))
	}
	if groupCount >= 20 { // 假设最大分组数为20
		return nil, errorx.New(errorx.ClientError, errorx.ErrGroupLimit, errorx.Message(errorx.ErrGroupLimit))
	}

	// 4. 创建分组
	gid := generateRandomString(8) // 生成12位随机字符串作为分组ID
	err = l.svcCtx.Query.TGroup.WithContext(l.ctx).Create(&model.TGroup{
		Gid:        gid,
		Username:   in.Username,
		Name:       in.GroupName,
		SortOrder:  0, // 初始排序号为0
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		DelFlag:    false,
	})
	if err != nil {
		return nil, errorx.New(errorx.SystemError, errorx.ErrInternalServer, errorx.Message(errorx.ErrInternalServer))
	}

	// 5. 添加到布隆过滤器
	err = l.svcCtx.BloomFilters.AddGroup(l.ctx, gid)
	if err != nil {
		logx.Errorf("添加分组到布隆过滤器失败: %v", err)
	}

	return &__.CommonResponse{
		Success: true,
		Message: "创建成功",
	}, nil
}
