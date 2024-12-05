package logic

import (
	"context"
	"shorterurl/user/rpc/internal/svc"
	"shorterurl/user/rpc/internal/types/errorx"
	__ "shorterurl/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserCheckUsernameLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserCheckUsernameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCheckUsernameLogic {
	return &UserCheckUsernameLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 检查用户名是否存在
func (l *UserCheckUsernameLogic) UserCheckUsername(in *__.CheckUsernameRequest) (*__.CheckUsernameResponse, error) {
	// 1. 使用布隆过滤器快速检查用户是否存在
	exists, err := l.svcCtx.BloomFilters.UserExists(l.ctx, in.Username)
	if err != nil {
		return nil, err
	}
	if !exists {
		return &__.CheckUsernameResponse{
			Exist: false,
		}, nil
	}

	// 2. 检查用户是否真实存在
	user, err := l.svcCtx.Query.TUser.WithContext(l.ctx).Where(l.svcCtx.Query.TUser.Username.Eq(in.Username)).First()
	if err != nil {
		return nil, errorx.New(errorx.SystemError, errorx.ErrInternalServer, "查询用户信息失败")
	}

	return &__.CheckUsernameResponse{
		Exist: user != nil,
	}, nil
}
