package logic

import (
	"context"

	"shorterurl/link/rpc/internal/svc"
	"shorterurl/link/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecycleBinPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRecycleBinPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecycleBinPageLogic {
	return &RecycleBinPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RecycleBinPageLogic) RecycleBinPage(in *pb.PageRecycleBinShortLinkRequest) (*pb.PageRecycleBinShortLinkResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.PageRecycleBinShortLinkResponse{}, nil
}
