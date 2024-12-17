package logic

import (
	"context"

	"shorterurl/link/rpc/internal/svc"
	"shorterurl/link/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecycleBinSaveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRecycleBinSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecycleBinSaveLogic {
	return &RecycleBinSaveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// --------------------- 回收站管理接口 ---------------------
func (l *RecycleBinSaveLogic) RecycleBinSave(in *pb.SaveToRecycleBinRequest) (*pb.SaveToRecycleBinResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.SaveToRecycleBinResponse{}, nil
}
