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

type CreateShortLinkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建短链接
func NewCreateShortLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateShortLinkLogic {
	return &CreateShortLinkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateShortLinkLogic) CreateShortLink(req *types.CreateLinkReq) (resp *types.CreateLinkResp, err error) {
	// 获取用户信息
	userInfo, ok := types.GetUserFromCtx(l.ctx)
	if !ok {
		return nil, errors.New("未找到用户信息")
	}

	// 构建RPC请求
	rpcReq := &shortlinkservice.CreateShortLinkRequest{
		OriginUrl:     req.OriginUrl,
		Gid:           req.Gid,
		ValidDateType: int32(req.ValidDateType),
		ValidDate:     req.ValidDate,
		Describe:      req.Describe,
		CreatedType:   int32(req.CreatedType),
	}

	// 添加元数据
	ctx := metadata.AppendToOutgoingContext(l.ctx, "username", userInfo.Username)
	rpcResp, err := l.svcCtx.LinkRpc.ShortLinkCreate(ctx, rpcReq)
	if err != nil {
		l.Logger.Errorf("创建短链接失败 username: %s, error: %v", userInfo.Username, err)
		return nil, err
	}

	l.Logger.Infof("创建短链接成功 username: %s, fullShortUrl: %s, describe: %s",
		userInfo.Username, rpcResp.FullShortUrl, req.Describe)

	// 构建响应
	return &types.CreateLinkResp{
		FullShortUrl: rpcResp.FullShortUrl,
		OriginUrl:    rpcResp.OriginUrl,
		Gid:          rpcResp.Gid,
	}, nil
}
