package logic

import (
	"context"
	"shorterurl/user/rpc/internal/constant"
	"shorterurl/user/rpc/internal/svc"
	"shorterurl/user/rpc/internal/types/errorx"
	__ "shorterurl/user/rpc/pb"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type GroupDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupDeleteLogic {
	return &GroupDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GroupDelete 删除分组
func (l *GroupDeleteLogic) GroupDelete(in *__.GroupDeleteRequest) (*__.CommonResponse, error) {
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

	// 创建分布式锁，防止并发删除
	lockKey := constant.LOCK_GROUP_DELETE_KEY + in.Gid
	lock := redis.NewRedisLock(l.svcCtx.Redis, lockKey)
	lock.SetExpire(30) // 设置锁的过期时间为30秒

	// 获取锁
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

	// 1. 执行软删除更新
	result, err := l.svcCtx.Query.TGroup.WithContext(l.ctx).
		Where(l.svcCtx.Query.TGroup.Gid.Eq(in.Gid)).
		Where(l.svcCtx.Query.TGroup.Username.Eq(username)).
		Where(l.svcCtx.Query.TGroup.DelFlag.Is(false)).
		UpdateSimple(
			l.svcCtx.Query.TGroup.DelFlag.Value(true),
			l.svcCtx.Query.TGroup.UpdateTime.Value(time.Now()),
		)
	if err != nil {
		return nil, errorx.New(errorx.SystemError, errorx.ErrInternalServer, errorx.Message(errorx.ErrInternalServer))
	}

	// 检查是否有记录被更新
	if result.RowsAffected == 0 {
		return nil, errorx.New(errorx.ClientError, errorx.ErrGroupNotFound, "分组不存在或无权限删除")
	}

	// 2. 调用短链接服务将该分组下的所有短链接移入回收站
	// TODO: 这里需要调用短链接服务的RPC接口
	// 示例：
	// err = l.svcCtx.LinkRpc.MoveGroupToRecycleBin(l.ctx, &linkservice.RecycleBinRequest{
	//     Username: username,
	//     Gid:      in.Gid,
	// })
	// if err != nil {
	//     logx.Errorf("移动短链接到回收站失败: username=%s, gid=%s, error=%v", username, in.Gid, err)
	//     // 注意：这里我们选择只记录日志而不返回错误，因为分组已经被删除
	// }

	return &__.CommonResponse{
		Success: true,
		Message: "删除成功",
	}, nil
}
