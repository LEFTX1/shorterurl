package group

import (
	"context"

	"go-zero-shorterurl/admin/internal/svc"
	"go-zero-shorterurl/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SortGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSortGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SortGroupLogic {
	return &SortGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SortGroupLogic) SortGroup(req *types.ShortLinkGroupSortReq) error {
	// todo: add your logic here and delete this line

	return nil
}
