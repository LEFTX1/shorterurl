package logic

import (
	"context"

	"shorterurl/link/rpc/internal/svc"
	"shorterurl/link/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type StatsAccessRecordQueryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStatsAccessRecordQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StatsAccessRecordQueryLogic {
	return &StatsAccessRecordQueryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *StatsAccessRecordQueryLogic) StatsAccessRecordQuery(in *pb.AccessRecordQueryRequest) (*pb.AccessRecordQueryResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.AccessRecordQueryResponse{}, nil
}
