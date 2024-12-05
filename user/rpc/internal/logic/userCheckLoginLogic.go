package logic

import (
	"context"
	"shorterurl/user/rpc/internal/constant"
	"shorterurl/user/rpc/internal/svc"
	"shorterurl/user/rpc/internal/types/errorx"
	__ "shorterurl/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserCheckLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserCheckLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCheckLoginLogic {
	return &UserCheckLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 检查用户登录状态
func (l *UserCheckLoginLogic) UserCheckLogin(in *__.CheckLoginRequest) (*__.CommonResponse, error) {
	// 1. 检查用户是否已登录
	loginKey := constant.UserLoginKey + in.Username
	exists, err := l.svcCtx.Redis.ExistsCtx(l.ctx, loginKey)
	if err != nil {
		return nil, errorx.New(errorx.SystemError, errorx.ErrInternalServer, "检查登录状态失败")
	}
	if !exists {
		return &__.CommonResponse{
			Success: false,
			Message: "用户未登录",
		}, nil
	}

	// 2. 更新token过期时间
	err = l.svcCtx.Redis.ExpireCtx(l.ctx, loginKey, 30*60) // 30分钟
	if err != nil {
		logx.Errorf("更新token过期时间失败: %v", err)
	}

	return &__.CommonResponse{
		Success: true,
		Message: "用户已登录",
	}, nil
}
