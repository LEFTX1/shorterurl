package logic

import (
	"context"

	"shorterurl/link/rpc/internal/svc"
	"shorterurl/link/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type StatsGetGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStatsGetGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StatsGetGroupLogic {
	return &StatsGetGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *StatsGetGroupLogic) StatsGetGroup(in *pb.GetGroupStatsRequest) (*pb.GetGroupStatsResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.GetGroupStatsResponse{}, nil
}
