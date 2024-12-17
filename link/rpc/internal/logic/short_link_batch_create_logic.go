package logic

import (
	"context"

	"shorterurl/link/rpc/internal/svc"
	"shorterurl/link/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShortLinkBatchCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewShortLinkBatchCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShortLinkBatchCreateLogic {
	return &ShortLinkBatchCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ShortLinkBatchCreateLogic) ShortLinkBatchCreate(in *pb.BatchCreateShortLinkRequest) (*pb.BatchCreateShortLinkResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.BatchCreateShortLinkResponse{}, nil
}
