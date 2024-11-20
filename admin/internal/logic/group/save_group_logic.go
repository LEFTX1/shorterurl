package group

import (
	"context"

	"go-zero-shorterurl/admin/internal/svc"
	"go-zero-shorterurl/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveGroupLogic {
	return &SaveGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveGroupLogic) SaveGroup(req *types.ShortLinkGroupSaveReq) error {
	// todo: add your logic here and delete this line

	return nil
}
