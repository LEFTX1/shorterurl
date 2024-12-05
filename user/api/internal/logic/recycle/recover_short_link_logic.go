package recycle

import (
	"context"
	"shorterurl/user/api/internal/svc"
	"shorterurl/user/api/internal/types"
	"shorterurl/user/api/internal/types/errorx"
	"shorterurl/user/rpc/userservice"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"
)

type RecoverShortLinkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRecoverShortLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecoverShortLinkLogic {
	return &RecoverShortLinkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RecoverShortLinkLogic) RecoverShortLink(req *types.RecycleBinRecoverReq) (resp *types.SuccessResp, err error) {
	// 从上下文中获取用户信息
	userInfo := l.ctx.Value(types.UserContextKey).(*types.UserInfo)
	if userInfo == nil {
		return nil, errorx.New(errorx.ClientError, errorx.ErrInternalServer, "未找到用户信息")
	}

	// 使用推荐的方式添加元数据
	ctx := metadata.AppendToOutgoingContext(l.ctx, "username", userInfo.Username)

	// 调用RPC服务恢复短链接
	_, err = l.svcCtx.UserRpc.RecycleBinRecover(ctx, &userservice.RecycleBinOperateRequest{
		Gid:          req.Gid,
		FullShortUrl: req.FullShortUrl,
	})
	if err != nil {
		logx.Errorf("恢复短链接失败 username: %s, gid: %s, url: %s, error: %v",
			userInfo.Username, req.Gid, req.FullShortUrl, err)
		return nil, err
	}

	return &types.SuccessResp{
		Code:    "0",
		Success: true,
	}, nil
}
