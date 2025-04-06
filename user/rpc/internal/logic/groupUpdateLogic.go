package logic

import (
	"context"
	"shorterurl/user/rpc/internal/constant"
	"shorterurl/user/rpc/internal/svc"
	"shorterurl/user/rpc/internal/types/errorx"
	__ "shorterurl/user/rpc/pb"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"google.golang.org/grpc/metadata"
)

type GroupUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupUpdateLogic {
	return &GroupUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GroupUpdate 更新分组
func (l *GroupUpdateLogic) GroupUpdate(in *__.GroupUpdateRequest) (*__.CommonResponse, error) {
	// 从metadata中获取用户名
	md, ok := metadata.FromIncomingContext(l.ctx)
	if !ok {
		return nil, errorx.New(errorx.ClientError, errorx.ErrInternalServer, "未找到用户信息")
	}
	usernames := md.Get("username")
	if len(usernames) == 0 {
		return nil, errorx.New(errorx.ClientError, errorx.ErrInternalServer, "未找到用户信息")
	}
	username := usernames[0]

	// 1. 检查分组是否存在（使用布隆过滤器快速判断）
	exists, err := l.svcCtx.BloomFilters.GroupExists(l.ctx, in.Gid)
	if err != nil {
		return nil, errorx.New(errorx.SystemError, errorx.ErrBloomFilterCheck, errorx.Message(errorx.ErrBloomFilterCheck))
	}
	if !exists {
		return nil, errorx.New(errorx.ClientError, errorx.ErrGroupNotFound, errorx.Message(errorx.ErrGroupNotFound))
	}

	// 2. 创建分布式锁，防止并发更新相同分组
	lockKey := constant.LockGroupUpdateKey + in.Gid
	lock := redis.NewRedisLock(l.svcCtx.Redis, lockKey)
	lock.SetExpire(30) // 设置锁的过期时间为30秒

	// 3. 获取锁
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

	// 4. 更新分组
	group, err := l.svcCtx.Query.TGroup.WithContext(l.ctx).
		Where(l.svcCtx.Query.TGroup.Username.Eq(username)).
		Where(l.svcCtx.Query.TGroup.Gid.Eq(in.Gid)).
		Where(l.svcCtx.Query.TGroup.DelFlag.Is(false)).
		First()
	if err != nil {
		return nil, errorx.New(errorx.SystemError, errorx.ErrInternalServer, errorx.Message(errorx.ErrInternalServer))
	}
	if group == nil {
		return nil, errorx.New(errorx.ClientError, errorx.ErrGroupNotFound, errorx.Message(errorx.ErrGroupNotFound))
	}

	// 验证用户权限
	if group.Username != username {
		return nil, errorx.New(errorx.ClientError, errorx.ErrGroupNotFound, "无权限更新该分组")
	}

	// 保存原始时间以便于日志记录
	originalTime := group.UpdateTime

	// 使用 username 和 gid 作为条件更新，而不是主键 ID
	currentTime := time.Now()

	// 确保新时间与原时间至少相差1秒
	if currentTime.Sub(originalTime) < time.Second {
		time.Sleep(time.Second - currentTime.Sub(originalTime))
		currentTime = time.Now()
	}

	l.Infof("更新分组，原时间: %v, 设置更新时间为: %v", originalTime, currentTime)

	_, err = l.svcCtx.Query.TGroup.WithContext(l.ctx).
		Where(l.svcCtx.Query.TGroup.Username.Eq(username)).
		Where(l.svcCtx.Query.TGroup.Gid.Eq(in.Gid)).
		Where(l.svcCtx.Query.TGroup.DelFlag.Is(false)).
		UpdateSimple(
			l.svcCtx.Query.TGroup.Name.Value(in.Name),
			l.svcCtx.Query.TGroup.UpdateTime.Value(currentTime),
		)
	if err != nil {
		return nil, errorx.New(errorx.SystemError, errorx.ErrInternalServer, errorx.Message(errorx.ErrInternalServer))
	}

	// 创建新的响应函数
	createSuccessResponse := func() *__.CommonResponse {
		return &__.CommonResponse{
			Success: true,
			Message: "更新成功",
		}
	}

	response := createSuccessResponse()

	// 额外检查
	if !response.GetSuccess() {
		l.Error("响应对象创建后Success字段为false，尝试手动修正")
		response.Success = true
	}

	l.Infof("返回响应: Success=%v, Message=%s", response.GetSuccess(), response.GetMessage())

	return response, nil
}
