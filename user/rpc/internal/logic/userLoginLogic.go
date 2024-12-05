package logic

import (
	"context"
	"fmt"
	"shorterurl/user/rpc/internal/constant"
	"shorterurl/user/rpc/internal/dal/query"
	"shorterurl/user/rpc/internal/svc"
	"shorterurl/user/rpc/internal/types/errorx"
	__ "shorterurl/user/rpc/pb"
	"time"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLoginLogic {
	return &UserLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户登录
func (l *UserLoginLogic) UserLogin(in *__.LoginRequest) (*__.LoginResponse, error) {
	// 1. 查询用户是否存在
	q := query.Use(l.svcCtx.DB)
	user, err := q.TUser.WithContext(l.ctx).Where(q.TUser.Username.Eq(in.Username), q.TUser.Password.Eq(in.Password)).First()
	if err != nil {
		return nil, errorx.New(errorx.ClientError, errorx.ErrUserNotFound, errorx.Message(errorx.ErrUserNotFound))
	}

	// 2. 检查用户是否已经登录
	loginKey := constant.UserLoginKey + in.Username
	existingToken, err := l.svcCtx.Redis.HgetallCtx(l.ctx, loginKey)
	if err == nil && len(existingToken) > 0 {
		// 如果已经登录，更新过期时间并返回已有token
		err = l.svcCtx.Redis.ExpireCtx(l.ctx, loginKey, int(30*time.Minute.Seconds()))
		if err != nil {
			logx.Errorf("更新登录token过期时间失败: %v", err)
		}
		// 返回第一个token
		for token := range existingToken {
			return &__.LoginResponse{
				Token:      token,
				Username:   user.Username,
				RealName:   user.RealName,
				CreateTime: time.Now().Format("2006-01-02 15:04:05"),
			}, nil
		}
	}

	// 3. 生成新的登录token
	token := uuid.New().String()

	// 4. 存储登录信息到Redis，只存储用户名作为标识
	err = l.svcCtx.Redis.HsetCtx(l.ctx, loginKey, token, user.Username)
	if err != nil {
		return nil, errorx.New(errorx.SystemError, errorx.ErrInternalServer, fmt.Sprintf("存储登录信息失败: %v", err))
	}
	err = l.svcCtx.Redis.ExpireCtx(l.ctx, loginKey, int(30*time.Minute.Seconds()))
	if err != nil {
		logx.Errorf("设置登录token过期时间失败: %v", err)
	}

	return &__.LoginResponse{
		Token:      token,
		Username:   user.Username,
		RealName:   user.RealName,
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}
