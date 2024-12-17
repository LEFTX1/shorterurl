package logic

import (
	"context"

	"shorterurl/link/rpc/internal/svc"
	"shorterurl/link/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShortLinkUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewShortLinkUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShortLinkUpdateLogic {
	return &ShortLinkUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ShortLinkUpdateLogic) ShortLinkUpdate(in *pb.UpdateShortLinkRequest) (*pb.UpdateShortLinkResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.UpdateShortLinkResponse{}, nil
}
