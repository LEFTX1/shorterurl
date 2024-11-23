package user

import (
	"context"
	"errors"
	"go-zero-shorterurl/admin/internal/dal/model"
	"go-zero-shorterurl/admin/internal/svc"
	"go-zero-shorterurl/admin/internal/types"
	"go-zero-shorterurl/admin/internal/types/errorx"
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

// user_info_logic.go
// user_info_logic.go
func (l *UserInfoLogic) UserInfo(req *types.UserUsernameReq) (resp *types.UserInfoResp, err error) {
	username := req.Username
	if username == "" {
		return nil, errorx.NewUserError("USERNAME_EMPTY")
	}

	// 使用 GORM 原生查询，避免使用 Gen 的查询表达式
	var user model.TUser
	err = l.svcCtx.DB.WithContext(l.ctx).
		Select("*").
		Where("username = ? AND del_flag = ?", username, false).
		First(&user).Error

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
		Phone:      user.Phone,
		Mail:       user.Mail,
		CreateTime: user.CreateTime.Format("2006-01-02 15:04:05"),
		UpdateTime: user.UpdateTime.Format("2006-01-02 15:04:05"),
	}

	return resp, nil
}
