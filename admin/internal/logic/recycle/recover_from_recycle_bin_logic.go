package recycle

import (
	"context"

	"go-zero-shorterurl/admin/internal/svc"
	"go-zero-shorterurl/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecoverFromRecycleBinLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRecoverFromRecycleBinLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecoverFromRecycleBinLogic {
	return &RecoverFromRecycleBinLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RecoverFromRecycleBinLogic) RecoverFromRecycleBin(req *types.RecycleBinRecoverReq) error {
	// todo: add your logic here and delete this line

	return nil
}
