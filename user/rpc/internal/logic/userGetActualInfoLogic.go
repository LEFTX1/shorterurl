package logic

import (
	"context"
	"errors"
	"shorterurl/user/rpc/internal/svc"
	"shorterurl/user/rpc/internal/types/errorx"
	__ "shorterurl/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UserGetActualInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserGetActualInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserGetActualInfoLogic {
	return &UserGetActualInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取用户信息（无脱敏）
func (l *UserGetActualInfoLogic) UserGetActualInfo(in *__.CheckUsernameRequest) (*__.UserInfoResponse, error) {
	// 1. 从数据库获取用户信息
	user, err := l.svcCtx.Query.TUser.WithContext(l.ctx).Where(l.svcCtx.Query.TUser.Username.Eq(in.Username)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.New(errorx.ClientError, errorx.ErrUserNotFound, "用户不存在")
		}
		return nil, errorx.New(errorx.SystemError, errorx.ErrInternalServer, "获取用户信息失败")
	}

	// 2. 返回用户信息（无需脱敏）
	return &__.UserInfoResponse{
		Username:   user.Username,
		RealName:   user.RealName,
		Phone:      user.Phone,
		Mail:       user.Mail,
		CreateTime: user.CreateTime.Format("2006-01-02 15:04:05"),
		UpdateTime: user.UpdateTime.Format("2006-01-02 15:04:05"),
	}, nil
}
