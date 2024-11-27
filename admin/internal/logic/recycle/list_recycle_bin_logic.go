package recycle

import (
	"context"

	"shorterurl/admin/internal/svc"
	"shorterurl/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRecycleBinLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListRecycleBinLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRecycleBinLogic {
	return &ListRecycleBinLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListRecycleBinLogic) ListRecycleBin(req *types.RecycleBinPageReq) (resp *types.RecycleBinResp, err error) {
	// todo: add your logic here and delete this line

	return
}
