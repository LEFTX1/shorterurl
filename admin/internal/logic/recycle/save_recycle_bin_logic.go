package recycle

import (
	"context"

	"go-zero-shorterurl/admin/internal/svc"
	"go-zero-shorterurl/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveRecycleBinLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveRecycleBinLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveRecycleBinLogic {
	return &SaveRecycleBinLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveRecycleBinLogic) SaveRecycleBin(req *types.RecycleBinSaveReq) error {
	// todo: add your logic here and delete this line

	return nil
}
