package link

import (
	"context"
	"errors"
	"shorterurl/link/rpc/shortlinkservice"
	"shorterurl/user/api/internal/svc"
	"shorterurl/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"
)

type UpdateShortLinkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新短链接
func NewUpdateShortLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateShortLinkLogic {
	return &UpdateShortLinkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateShortLinkLogic) UpdateShortLink(req *types.UpdateLinkReq) (resp *types.SuccessResp, err error) {
	// 获取用户信息
	userInfo, ok := types.GetUserFromCtx(l.ctx)
	if !ok {
		return nil, errors.New("未找到用户信息")
	}

	// 构建RPC请求
	rpcReq := &shortlinkservice.UpdateShortLinkRequest{
		OriginUrl:     req.OriginUrl,
		FullShortUrl:  req.FullShortUrl,
		Gid:           req.Gid,
		ValidDateType: int32(req.ValidDateType),
		ValidDate:     req.ValidDate,
		Describe:      req.Describe,
	}

	// 添加元数据
	ctx := metadata.AppendToOutgoingContext(l.ctx, "username", userInfo.Username)
	_, err = l.svcCtx.LinkRpc.ShortLinkUpdate(ctx, rpcReq)
	if err != nil {
		l.Logger.Errorf("更新短链接失败 username: %s, fullShortUrl: %s, error: %v",
			userInfo.Username, req.FullShortUrl, err)
		return nil, err
	}

	// 构建响应
	return &types.SuccessResp{
		Code:    "0",
		Success: true,
	}, nil
}
