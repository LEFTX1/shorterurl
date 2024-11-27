package user

import (
	"context"
	"shorterurl/admin/internal/svc"
	"shorterurl/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// UserRegisterLogic 处理用户注册逻辑的结构体
type UserRegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUserRegisterLogic 创建一个新的 UserRegisterLogic 实例
func NewUserRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRegisterLogic {
	return &UserRegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UserRegister 处理用户注册请求
func (l *UserRegisterLogic) UserRegister(req *types.UserRegisterReq) (*types.UserRegisterResp, error) {

}
