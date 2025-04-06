package logic

import (
	"context"
	"errors"
	"shorterurl/user/rpc/internal/svc"
	"shorterurl/user/rpc/internal/types/errorx"
	__ "shorterurl/user/rpc/pb"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UserGetInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserGetInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserGetInfoLogic {
	return &UserGetInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取用户信息（带脱敏）
func (l *UserGetInfoLogic) UserGetInfo(in *__.CheckUsernameRequest) (*__.UserInfoResponse, error) {
	// 1. 从数据库获取用户信息
	user, err := l.svcCtx.Query.TUser.WithContext(l.ctx).Where(l.svcCtx.Query.TUser.Username.Eq(in.Username)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.New(errorx.ClientError, errorx.ErrUserNotFound, "用户不存在")
		}
		return nil, errorx.New(errorx.SystemError, errorx.ErrInternalServer, "获取用户信息失败")
	}

	// 2. 脱敏处理
	phone := user.Phone
	if len(phone) > 7 {
		phone = phone[:3] + "****" + phone[len(phone)-4:]
	}

	mail := user.Mail
	if len(mail) > 7 {
		atIndex := strings.Index(mail, "@")
		if atIndex > 0 {
			prefix := mail[:atIndex]
			if len(prefix) > 3 {
				prefix = prefix[:3] + "****"
			}
			mail = prefix + mail[atIndex:]
		}
	}

	// 3. 返回用户信息
	return &__.UserInfoResponse{
		Id:         user.ID,
		Username:   user.Username,
		RealName:   user.RealName,
		Phone:      phone,
		Mail:       mail,
		CreateTime: user.CreateTime.Format("2006-01-02 15:04:05"),
		UpdateTime: user.UpdateTime.Format("2006-01-02 15:04:05"),
	}, nil
}
