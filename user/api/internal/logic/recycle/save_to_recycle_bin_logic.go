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

type SaveToRecycleBinLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveToRecycleBinLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveToRecycleBinLogic {
	return &SaveToRecycleBinLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveToRecycleBinLogic) SaveToRecycleBin(req *types.RecycleBinSaveReq) (resp *types.SuccessResp, err error) {
	// 从上下文中获取用户信息
	userInfo := l.ctx.Value(types.UserContextKey).(*types.UserInfo)
	if userInfo == nil {
		return nil, errorx.New(errorx.ClientError, errorx.ErrInternalServer, "未找到用户信息")
	}

	// 使用推荐的方式添加元数据
	ctx := metadata.AppendToOutgoingContext(l.ctx, "username", userInfo.Username)

	// 调用RPC服务保存到回收站
	_, err = l.svcCtx.UserRpc.RecycleBinSave(ctx, &userservice.RecycleBinOperateRequest{
		Gid:          req.Gid,
		FullShortUrl: req.FullShortUrl,
	})
	if err != nil {
		logx.Errorf("保存到回收站失败 username: %s, gid: %s, url: %s, error: %v",
			userInfo.Username, req.Gid, req.FullShortUrl, err)
		return nil, err
	}

	return &types.SuccessResp{
		Code:    "0",
		Success: true,
	}, nil
}
