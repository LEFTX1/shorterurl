// Package user internal/logic/user/check_username_logic.go
package user

import (
	"context"

	"go-zero-shorterurl/admin/internal/svc"
	"go-zero-shorterurl/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckUsernameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckUsernameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckUsernameLogic {
	return &CheckUsernameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckUsernameLogic) CheckUsername(req *types.UserCheckUsernameReq) (resp bool, err error) {
	// 1 请求中获取用户名
	userName := req.Username

	// 2 从布隆过滤器里判断用户名是否存在
	resp, err = l.svcCtx.BloomFilters.UserExists(l.ctx, userName)

	// 2.1 处理错误
	if err != nil {
		l.Logger.Errorf("判断用户名是否存在失败: %v, username: %s", err, userName)
		return false, err
	}

	// 3 返回结果
	return resp, nil
}
