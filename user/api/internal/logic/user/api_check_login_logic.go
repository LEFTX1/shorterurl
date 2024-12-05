package user

import (
	"context"
	"shorterurl/user/api/internal/svc"
	"shorterurl/user/api/internal/types"
	"shorterurl/user/api/internal/types/errorx"
	"shorterurl/user/rpc/userservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiCheckLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiCheckLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiCheckLoginLogic {
	return &ApiCheckLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiCheckLoginLogic) ApiCheckLogin(req *types.UserCheckLoginReq) (resp *types.SuccessResp, err error) {
	// 参数校验
	if len(req.Username) == 0 || len(req.Token) == 0 {
		return nil, errorx.New(errorx.ClientError, "INVALID_PARAMS", "用户名或令牌不能为空")
	}

	// 调用rpc服务检查登录状态
	rpcResp, err := l.svcCtx.UserRpc.UserCheckLogin(l.ctx, &userservice.CheckLoginRequest{
		Username: req.Username,
		Token:    req.Token,
	})
	if err != nil {
		logx.Errorf("检查用户登录状态RPC调用失败 username: %s, error: %v", req.Username, err)
		return nil, errorx.New(errorx.SystemError, "RPC_ERROR", "检查登录状态失败")
	}

	return &types.SuccessResp{
		Success: rpcResp.GetSuccess(),
		Code:    rpcResp.GetMessage(),
	}, nil
}
