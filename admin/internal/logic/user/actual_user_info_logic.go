package user

import (
	"context"
	"errors"

	"go-zero-shorterurl/admin/internal/dal/model"
	"go-zero-shorterurl/admin/internal/svc"
	"go-zero-shorterurl/admin/internal/types"
	"go-zero-shorterurl/admin/internal/types/errorx"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type ActualUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewActualUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ActualUserInfoLogic {
	return &ActualUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ActualUserInfoLogic) ActualUserInfo(req *types.UserUsernameReq) (resp *types.UserInfoResp, err error) {
	username := req.Username
	if username == "" {
		return nil, errorx.NewUserError("USERNAME_EMPTY")
	}

	// 查询用户信息（获取未脱敏的实际数据）
	var user *model.TUser
	user, err = l.svcCtx.Query.TUser.WithContext(l.ctx).
		Where(l.svcCtx.Query.TUser.Username.Eq(username)).
		Where(l.svcCtx.Query.TUser.DelFlag.Is(false)).
		First()

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.NewUserError("USER_NOT_EXIST")
		}
		logx.WithContext(l.ctx).Errorf("查询实际用户信息失败: %v, username: %s", err, username)
		return nil, errorx.NewSystemError("DB_ERROR")
	}

	// 返回未脱敏的用户信息
	resp = &types.UserInfoResp{
		Username:   user.Username,
		RealName:   user.RealName,
		Phone:      user.Phone, // 返回原始手机号
		Mail:       user.Mail,  // 返回原始邮箱
		CreateTime: user.CreateTime.Format("2006-01-02 15:04:05"),
		UpdateTime: user.UpdateTime.Format("2006-01-02 15:04:05"),
	}

	return resp, nil
}
