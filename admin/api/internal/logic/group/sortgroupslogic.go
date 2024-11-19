package group

import (
	"context"

	"go-zero-shorterurl/admin/api/internal/svc"
	"go-zero-shorterurl/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SortGroupsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSortGroupsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SortGroupsLogic {
	return &SortGroupsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SortGroupsLogic) SortGroups(req *types.GroupSortReq) error {
	// todo: add your logic here and delete this line

	return nil
}
