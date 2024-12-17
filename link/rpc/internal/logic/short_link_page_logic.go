package logic

import (
	"context"

	"shorterurl/link/rpc/internal/svc"
	"shorterurl/link/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShortLinkPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewShortLinkPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShortLinkPageLogic {
	return &ShortLinkPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ShortLinkPageLogic) ShortLinkPage(in *pb.PageShortLinkRequest) (*pb.PageShortLinkResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.PageShortLinkResponse{}, nil
}
