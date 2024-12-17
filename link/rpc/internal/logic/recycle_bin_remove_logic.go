package logic

import (
	"context"

	"shorterurl/link/rpc/internal/svc"
	"shorterurl/link/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecycleBinRemoveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRecycleBinRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecycleBinRemoveLogic {
	return &RecycleBinRemoveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RecycleBinRemoveLogic) RecycleBinRemove(in *pb.RemoveFromRecycleBinRequest) (*pb.RemoveFromRecycleBinResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.RemoveFromRecycleBinResponse{}, nil
}
