package logic

import (
	"context"
	"shorterurl/user/rpc/internal/svc"
	"shorterurl/user/rpc/internal/types/errorx"
	__ "shorterurl/user/rpc/pb"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type UserUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserUpdateLogic {
	return &UserUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新用户信息
func (l *UserUpdateLogic) UserUpdate(in *__.UpdateRequest) (*__.CommonResponse, error) {
	// 1. 使用布隆过滤器快速检查用户是否存在
	exists, err := l.svcCtx.BloomFilters.UserExists(l.ctx, in.Username)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errorx.New(errorx.ClientError, errorx.ErrUserNotFound, "用户不存在")
	}

	// 2. 检查用户是否真实存在
	user, err := l.svcCtx.Query.TUser.WithContext(l.ctx).Where(l.svcCtx.Query.TUser.Username.Eq(in.Username)).First()
	if err != nil {
		return nil, errorx.New(errorx.SystemError, errorx.ErrInternalServer, "查询用户信息失败")
	}
	if user == nil {
		return nil, errorx.New(errorx.ClientError, errorx.ErrUserNotFound, "用户不存在")
	}

	// 3. 准备更新数据
	updates := make(map[string]interface{})
	if in.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errorx.New(errorx.SystemError, errorx.ErrInternalServer, "密码加密失败")
		}
		updates["password"] = string(hashedPassword)
	}
	if in.RealName != "" {
		updates["real_name"] = in.RealName
	}
	if in.Phone != "" {
		updates["phone"] = in.Phone
	}
	if in.Mail != "" {
		updates["mail"] = in.Mail
	}
	if len(updates) > 0 {
		updates["update_time"] = time.Now()
	}

	// 4. 执行更新
	if len(updates) > 0 {
		result, err := l.svcCtx.Query.TUser.WithContext(l.ctx).
			Where(l.svcCtx.Query.TUser.Username.Eq(in.Username)).
			Updates(updates)
		if err != nil {
			return nil, errorx.New(errorx.SystemError, errorx.ErrInternalServer, "更新用户信息失败")
		}
		if result.RowsAffected == 0 {
			return nil, errorx.New(errorx.SystemError, errorx.ErrInternalServer, "更新用户信息失败")
		}
	}

	return &__.CommonResponse{
		Success: true,
		Message: "更新成功",
	}, nil
}
