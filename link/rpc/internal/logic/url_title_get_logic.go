package logic

import (
	"context"

	"shorterurl/link/rpc/internal/svc"
	"shorterurl/link/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UrlTitleGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUrlTitleGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UrlTitleGetLogic {
	return &UrlTitleGetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// --------------------- URL标题功能接口 ---------------------
func (l *UrlTitleGetLogic) UrlTitleGet(in *pb.GetUrlTitleRequest) (*pb.GetUrlTitleResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.GetUrlTitleResponse{}, nil
}
