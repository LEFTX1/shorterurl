package logic

import (
	"context"
	"shorterurl/user/rpc/internal/constant"
	"shorterurl/user/rpc/internal/svc"
	"shorterurl/user/rpc/internal/types/errorx"
	__ "shorterurl/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserLogoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLogoutLogic {
	return &UserLogoutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户退出登录
func (l *UserLogoutLogic) UserLogout(in *__.LogoutRequest) (*__.CommonResponse, error) {
	// 1. 检查用户是否已登录
	loginKey := constant.UserLoginKey + in.Username
	exists, err := l.svcCtx.Redis.ExistsCtx(l.ctx, loginKey)
	if err != nil {
		return nil, errorx.New(errorx.SystemError, errorx.ErrInternalServer, "检查登录状态失败")
	}
	if !exists {
		return nil, errorx.New(errorx.ClientError, errorx.ErrUserNotFound, "用户未登录")
	}

	// 2. 验证token是否匹配
	_, err = l.svcCtx.Redis.HgetCtx(l.ctx, loginKey, in.Token)
	if err != nil {
		return nil, errorx.New(errorx.ClientError, errorx.ErrUserNotFound, "登录token无效")
	}

	// 3. 删除登录信息
	val, err := l.svcCtx.Redis.DelCtx(l.ctx, loginKey)
	if err != nil || val == 0 {
		logx.Errorf("删除登录信息失败: %v", err)
		return nil, errorx.New(errorx.SystemError, errorx.ErrInternalServer, "退出登录失败")
	}

	return &__.CommonResponse{
		Success: true,
		Message: "退出登录成功",
	}, nil
}
