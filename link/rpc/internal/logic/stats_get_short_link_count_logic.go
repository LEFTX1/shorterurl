package logic

import (
	"context"

	"shorterurl/link/rpc/internal/svc"
	"shorterurl/link/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type StatsGetShortLinkCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStatsGetShortLinkCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StatsGetShortLinkCountLogic {
	return &StatsGetShortLinkCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *StatsGetShortLinkCountLogic) StatsGetShortLinkCount(in *pb.GetShortLinkCountRequest) (*pb.GetShortLinkCountResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.GetShortLinkCountResponse{}, nil
}
