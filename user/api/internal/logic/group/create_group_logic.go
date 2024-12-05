package group

import (
	"context"
	"shorterurl/user/api/internal/svc"
	"shorterurl/user/api/internal/types"
	"shorterurl/user/api/internal/types/errorx"
	"shorterurl/user/rpc/userservice"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"
)

type CreateGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupLogic {
	return &CreateGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateGroupLogic) CreateGroup(req *types.ShortLinkGroupSaveReq) (resp *types.SuccessResp, err error) {
	// 从上下文中获取用户信息
	userInfo := l.ctx.Value(types.UserContextKey).(*types.UserInfo)
	if userInfo == nil {
		return nil, errorx.New(errorx.ClientError, errorx.ErrInternalServer, "未找到用户信息")
	}

	// 使用推荐的方式添加元数据
	ctx := metadata.AppendToOutgoingContext(l.ctx, "username", userInfo.Username)

	// 调用RPC服务创建分组
	_, err = l.svcCtx.UserRpc.GroupCreate(ctx, &userservice.GroupSaveRequest{
		Username:  userInfo.Username,
		GroupName: req.Name,
	})
	if err != nil {
		logx.Errorf("创建分组失败 username: %s, groupName: %s, error: %v", userInfo.Username, req.Name, err)
		return nil, err
	}

	return &types.SuccessResp{
		Code:    "0",
		Success: true,
	}, nil
}
