package user

import (
	"context"
	"shorterurl/user/api/internal/svc"
	"shorterurl/user/api/internal/types"
	"shorterurl/user/rpc/userservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiLogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiLogoutLogic {
	return &ApiLogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiLogoutLogic) ApiLogout(req *types.UserLogOutReq) (resp *types.SuccessResp, err error) {
	// 调用 RPC 服务处理登出请求
	rpcResp, err := l.svcCtx.UserRpc.UserLogout(l.ctx, &userservice.LogoutRequest{
		Username: req.Username,
		Token:    req.Token,
	})
	if err != nil {
		logx.Errorf("用户登出失败 username: %s, error: %v", req.Username, err)
		return nil, err
	}

	// 转换 RPC 响应为 API 响应
	return &types.SuccessResp{
		Code:    "200", // 成功状态码
		Success: rpcResp.Success,
	}, nil
}
