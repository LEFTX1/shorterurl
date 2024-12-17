package logic

import (
	"context"

	"shorterurl/link/rpc/internal/svc"
	"shorterurl/link/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecycleBinRecoverLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRecycleBinRecoverLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecycleBinRecoverLogic {
	return &RecycleBinRecoverLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RecycleBinRecoverLogic) RecycleBinRecover(in *pb.RecoverFromRecycleBinRequest) (*pb.RecoverFromRecycleBinResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.RecoverFromRecycleBinResponse{}, nil
}
