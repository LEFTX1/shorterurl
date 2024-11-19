package recycle

import (
	"context"

	"go-zero-shorterurl/admin/api/internal/svc"
	"go-zero-shorterurl/admin/api/internal/types"

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

func (l *ListRecycleBinLogic) ListRecycleBin(req *types.RecycleBinPageReq) (resp *types.ShortLinkPageResp, err error) {
	// todo: add your logic here and delete this line

	return
}
