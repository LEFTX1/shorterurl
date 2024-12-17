package logic

import (
	"context"

	"shorterurl/link/rpc/internal/svc"
	"shorterurl/link/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type StatsGroupAccessRecordQueryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStatsGroupAccessRecordQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StatsGroupAccessRecordQueryLogic {
	return &StatsGroupAccessRecordQueryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *StatsGroupAccessRecordQueryLogic) StatsGroupAccessRecordQuery(in *pb.GroupAccessRecordQueryRequest) (*pb.GroupAccessRecordQueryResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.GroupAccessRecordQueryResponse{}, nil
}
