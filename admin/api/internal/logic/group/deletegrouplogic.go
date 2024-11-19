package group

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"go-zero-shorterurl/admin/api/internal/svc"
)

type DeleteGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteGroupLogic {
	return &DeleteGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteGroupLogic) DeleteGroup() error {
	// todo: add your logic here and delete this line

	return nil
}
