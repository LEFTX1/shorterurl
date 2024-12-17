package logic

import (
	"context"

	"shorterurl/link/rpc/internal/svc"
	"shorterurl/link/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type StatsGetSingleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStatsGetSingleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StatsGetSingleLogic {
	return &StatsGetSingleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// --------------------- 短链接统计接口 ---------------------
func (l *StatsGetSingleLogic) StatsGetSingle(in *pb.GetSingleStatsRequest) (*pb.GetSingleStatsResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.GetSingleStatsResponse{}, nil
}
