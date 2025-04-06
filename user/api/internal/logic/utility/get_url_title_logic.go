package utility

import (
	"context"
	"shorterurl/link/rpc/shortlinkservice"
	"shorterurl/user/api/internal/svc"
	"shorterurl/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUrlTitleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取网站标题
func NewGetUrlTitleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUrlTitleLogic {
	return &GetUrlTitleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUrlTitleLogic) GetUrlTitle(req *types.GetUrlTitleReq) (resp *types.UrlTitleResp, err error) {
	// 调用link RPC服务获取URL标题
	result, err := l.svcCtx.LinkRpc.UrlTitleGet(l.ctx, &shortlinkservice.GetUrlTitleRequest{
		Url: req.Url,
	})
	if err != nil {
		logx.Errorf("获取URL标题失败: %v", err)
		return nil, err
	}

	// 构建响应
	resp = &types.UrlTitleResp{
		Data: result.Title,
	}

	return resp, nil
}
