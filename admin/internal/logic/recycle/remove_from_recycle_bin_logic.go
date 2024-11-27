package recycle

import (
	"context"

	"shorterurl/admin/internal/svc"
	"shorterurl/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveFromRecycleBinLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemoveFromRecycleBinLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveFromRecycleBinLogic {
	return &RemoveFromRecycleBinLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveFromRecycleBinLogic) RemoveFromRecycleBin(req *types.RecycleBinRemoveReq) error {
	// todo: add your logic here and delete this line

	return nil
}
