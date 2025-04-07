package logic

import (
	"context"
	"io"
	"shorterurl/user/rpc/internal/constant"
	"shorterurl/user/rpc/internal/dal/query"
	"shorterurl/user/rpc/internal/svc"
	"shorterurl/user/rpc/internal/types/errorx"
	__ "shorterurl/user/rpc/pb"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type GroupSortLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupSortLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupSortLogic {
	return &GroupSortLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GroupSort 分组排序
func (l *GroupSortLogic) GroupSort(stream __.UserService_GroupSortServer) error {
	// 从metadata中获取用户名
	md, ok := metadata.FromIncomingContext(l.ctx)
	if !ok {
		return errorx.New(errorx.ClientError, errorx.ErrInternalServer, "未找到用户信息")
	}
	usernames := md.Get("username")
	if len(usernames) == 0 {
		return errorx.New(errorx.ClientError, errorx.ErrInternalServer, "未找到用户信息")
	}
	username := usernames[0]

	// 创建分布式锁，防止并发排序
	lockKey := constant.LockGroupSortKey + username
	lock := redis.NewRedisLock(l.svcCtx.Redis, lockKey)
	lock.SetExpire(30) // 设置锁的过期时间为30秒

	// 获取锁
	acquired, err := lock.AcquireCtx(l.ctx)
	if err != nil {
		return errorx.New(errorx.SystemError, errorx.ErrDistributedLock, errorx.Message(errorx.ErrDistributedLock))
	}
	if !acquired {
		return errorx.New(errorx.ClientError, errorx.ErrTooManyRequests, errorx.Message(errorx.ErrTooManyRequests))
	}

	// 确保在函数结束时释放锁
	defer func() {
		if released, err := lock.Release(); err != nil {
			logx.Errorf("释放锁失败: %v", err)
		} else if !released {
			logx.Error("锁未被主动释放")
		}
	}()

	// 接收所有排序请求
	var sortRequests []*__.GroupSortRequest
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errorx.New(errorx.SystemError, errorx.ErrInternalServer, "接收排序请求失败")
		}
		sortRequests = append(sortRequests, req)
	}

	// 开启事务处理所有排序请求
	err = l.svcCtx.Query.Transaction(func(tx *query.Query) error {
		now := time.Now()
		for _, req := range sortRequests {
			// 检查分组是否存在且属于当前用户
			group, err := tx.TGroup.WithContext(l.ctx).
				Where(tx.TGroup.Gid.Eq(req.Gid)).
				Where(tx.TGroup.Username.Eq(username)).
				Where(tx.TGroup.DelFlag.Is(false)).
				First()
			if err != nil {
				return errorx.New(errorx.SystemError, errorx.ErrInternalServer, "查询分组失败")
			}
			if group == nil {
				return errorx.New(errorx.ClientError, errorx.ErrGroupNotFound, "分组不存在")
			}

			// 更新排序号
			_, err = tx.TGroup.WithContext(l.ctx).
				Where(tx.TGroup.ID.Eq(group.ID)).
				Where(tx.TGroup.Username.Eq(username)).
				UpdateSimple(
					tx.TGroup.SortOrder.Value(req.SortOrder),
					tx.TGroup.UpdateTime.Value(now),
				)
			if err != nil {
				return errorx.New(errorx.SystemError, errorx.ErrInternalServer, "更新排序失败")
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	// 发送响应
	return stream.SendAndClose(&__.CommonResponse{
		Success: true,
		Message: "排序成功",
	})
}
