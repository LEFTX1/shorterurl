// internal/logic/user/user_login_logic.go
package user

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go-zero-shorterurl/admin/internal/dal/model"
	"go-zero-shorterurl/admin/internal/svc"
	"go-zero-shorterurl/admin/internal/types"
	"go-zero-shorterurl/admin/internal/types/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

const (
	// 用户登录信息过期时间（30分钟）
	loginExpireTime = int(30 * time.Minute / time.Second)
)

func NewUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLoginLogic {
	return &UserLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
func (l *UserLoginLogic) UserLogin(req *types.UserLoginReq) (*types.UserLoginResp, error) {
	// 1. 先通过布隆过滤器快速判断用户是否存在
	exists, err := l.svcCtx.BloomFilters.UserExists(l.ctx, req.Username)
	if err != nil {
		return nil, errorx.NewSystemError(errorx.ServiceError)
	}
	if !exists {
		// 用户名不存在，直接返回错误
		return nil, errorx.NewUserError(errorx.UserNotExistError)
	}

	// 2. 查询用户是否存在并验证密码
	var user model.TUser
	err = l.svcCtx.DB.Where("username = ? AND password = ? AND del_flag = ?",
		req.Username, req.Password, false).First(&user).Error
	if err != nil {
		return nil, errorx.NewUserError(errorx.UserNotExistError)
	}

	// 3. 检查是否已经登录
	loginKey := fmt.Sprintf("login:%s", req.Username)
	result, err := l.svcCtx.Redis.HgetallCtx(l.ctx, loginKey)
	if err != nil {
		return nil, errorx.NewSystemError(errorx.ServiceError)
	}

	// 4. 如果已经登录，直接返回已有token
	if len(result) > 0 {
		err = l.svcCtx.Redis.ExpireCtx(l.ctx, loginKey, loginExpireTime)
		if err != nil {
			logx.Errorf("刷新登录过期时间失败: %v", err)
		}

		for token := range result {
			return &types.UserLoginResp{
				Token: token,
			}, nil
		}
	}

	// 5. 生成新的登录token
	token := uuid.New().String()

	// 6. 将用户信息序列化
	userJson, err := json.Marshal(user)
	if err != nil {
		return nil, errorx.NewSystemError(errorx.ServiceError)
	}

	// 7. 存储登录信息到Redis
	err = l.svcCtx.Redis.HsetCtx(l.ctx, loginKey, token, string(userJson))
	if err != nil {
		return nil, errorx.NewSystemError(errorx.ServiceError)
	}

	// 8. 设置过期时间
	err = l.svcCtx.Redis.ExpireCtx(l.ctx, loginKey, loginExpireTime)
	if err != nil {
		logx.Errorf("设置登录过期时间失败: %v", err)
	}

	// 9. 返回登录成功响应
	return &types.UserLoginResp{
		Token: token,
	}, nil
}
