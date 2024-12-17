package logic

import (
	"context"

	"shorterurl/link/rpc/internal/svc"
	"shorterurl/link/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShortLinkCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewShortLinkCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShortLinkCreateLogic {
	return &ShortLinkCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// --------------------- 短链接管理接口 ---------------------
func (l *ShortLinkCreateLogic) ShortLinkCreate(in *pb.CreateShortLinkRequest) (*pb.CreateShortLinkResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.CreateShortLinkResponse{}, nil
}
