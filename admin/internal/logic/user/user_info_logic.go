package user

import (
	"context"
	"errors"
	"shorterurl/admin/internal/common"
	"shorterurl/admin/internal/dal/model"
	"shorterurl/admin/internal/svc"
	"shorterurl/admin/internal/types"
	"shorterurl/admin/internal/types/errorx"

	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UserInfo user_info_logic.go
func (l *UserInfoLogic) UserInfo(req *types.UserUsernameReq) (resp *types.UserInfoResp, err error) {
	username := req.Username
	if username == "" {
		return nil, errorx.NewUserError("USERNAME_EMPTY")
	}

	// 使用gen
	var user *model.TUser
	user, err = l.svcCtx.Query.TUser.WithContext(l.ctx).
		Where(l.svcCtx.Query.TUser.Username.Eq(username)).
		Where(l.svcCtx.Query.TUser.DelFlag.Is(false)).
		First()

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.NewUserError("USER_NOT_EXIST")
		}
		logx.WithContext(l.ctx).Errorf("查询用户信息失败: %v, username: %s", err, username)
		return nil, errorx.NewSystemError("DB_ERROR")
	}

	resp = &types.UserInfoResp{
		Username:   user.Username,
		RealName:   user.RealName,
		Phone:      common.MaskPhone(user.Phone),
		Mail:       common.MaskEmail(user.Mail),
		CreateTime: user.CreateTime.Format("2006-01-02 15:04:05"),
		UpdateTime: user.UpdateTime.Format("2006-01-02 15:04:05"),
	}

	return resp, nil
}
