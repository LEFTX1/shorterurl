package group

import (
	"context"
	"io"
	"shorterurl/user/api/internal/svc"
	"shorterurl/user/api/internal/types"
	"shorterurl/user/api/internal/types/errorx"
	"shorterurl/user/rpc/userservice"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"
)

type ListGroupsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListGroupsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListGroupsLogic {
	return &ListGroupsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListGroupsLogic) ListGroups() (resp []types.ShortLinkGroupResp, err error) {
	// 从上下文中获取用户信息
	userInfo := l.ctx.Value(types.UserContextKey).(*types.UserInfo)
	if userInfo == nil {
		return nil, errorx.New(errorx.ClientError, errorx.ErrInternalServer, "未找到用户信息")
	}

	// 使用推荐的方式添加元数据
	ctx := metadata.AppendToOutgoingContext(l.ctx, "username", userInfo.Username)

	// 调用RPC服务获取分组列表
	stream, err := l.svcCtx.UserRpc.GroupList(ctx, &userservice.CommonRequest{})
	if err != nil {
		logx.Errorf("获取分组列表失败 username: %s, error: %v", userInfo.Username, err)
		return nil, err
	}

	// 处理流式响应
	var groups []types.ShortLinkGroupResp
	for {
		group, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			logx.Errorf("接收分组数据失败 username: %s, error: %v", userInfo.Username, err)
			return nil, err
		}
		groups = append(groups, types.ShortLinkGroupResp{
			Gid:            group.Gid,
			Name:           group.Name,
			SortOrder:      int(group.SortOrder),
			ShortLinkCount: int(group.ShortLinkCount),
		})
	}

	return groups, nil
}
