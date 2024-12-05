package logic

import (
	"context"
	"shorterurl/user/rpc/internal/svc"
	"shorterurl/user/rpc/internal/types/errorx"
	__ "shorterurl/user/rpc/pb"

	"google.golang.org/grpc/metadata"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupListLogic {
	return &GroupListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GroupList 获取分组列表
func (l *GroupListLogic) GroupList(in *__.CommonRequest, stream __.UserService_GroupListServer) error {
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

	// 查询用户的所有分组，按排序号和更新时间降序排列
	groups, err := l.svcCtx.Query.TGroup.WithContext(l.ctx).
		Where(l.svcCtx.Query.TGroup.Username.Eq(username)).
		Where(l.svcCtx.Query.TGroup.DelFlag.Is(false)).
		Order(l.svcCtx.Query.TGroup.SortOrder.Desc(), l.svcCtx.Query.TGroup.UpdateTime.Desc()).
		Find()
	if err != nil {
		return errorx.New(errorx.SystemError, errorx.ErrInternalServer, "获取分组列表失败")
	}

	// 收集所有分组的 GID
	var gids []string
	for _, group := range groups {
		gids = append(gids, group.Gid)
	}

	// 批量查询所有分组的短链接数量
	//TODO 未来这里调用link的rpc服务

	return nil
}
