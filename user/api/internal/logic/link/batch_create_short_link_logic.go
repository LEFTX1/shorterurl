package link

import (
	"context"
	"errors"
	"shorterurl/link/rpc/shortlinkservice"
	"shorterurl/user/api/internal/svc"
	"shorterurl/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"
)

type BatchCreateShortLinkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量创建短链接
func NewBatchCreateShortLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchCreateShortLinkLogic {
	return &BatchCreateShortLinkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchCreateShortLinkLogic) BatchCreateShortLink(req *types.BatchCreateLinkReq) (resp *types.BatchCreateLinkResp, err error) {
	// 获取用户信息
	userInfo, ok := types.GetUserFromCtx(l.ctx)
	if !ok {
		return nil, errors.New("未找到用户信息")
	}

	// 验证URL列表非空
	if len(req.OriginUrls) == 0 {
		return nil, errors.New("原始URL列表不能为空")
	}

	// 构建RPC请求
	rpcReq := &shortlinkservice.BatchCreateShortLinkRequest{
		OriginUrls:    req.OriginUrls,
		Gid:           req.Gid,
		ValidDateType: int32(req.ValidDateType),
		ValidDate:     req.ValidDate,
		Describe:      req.Describes[0], // 批量创建时使用第一个描述
	}

	// 添加元数据
	ctx := metadata.AppendToOutgoingContext(l.ctx, "username", userInfo.Username)
	rpcResp, err := l.svcCtx.LinkRpc.ShortLinkBatchCreate(ctx, rpcReq)
	if err != nil {
		l.Logger.Errorf("批量创建短链接失败 username: %s, error: %v", userInfo.Username, err)
		return nil, err
	}

	// 构建响应
	linkBaseInfos := make([]types.LinkBaseInfo, 0, len(rpcResp.Results))
	for _, result := range rpcResp.Results {
		linkBaseInfos = append(linkBaseInfos, types.LinkBaseInfo{
			FullShortUrl: result.FullShortUrl,
			OriginUrl:    result.OriginUrl,
			Describe:     req.Describes[0], // 使用请求中的描述
		})
	}

	return &types.BatchCreateLinkResp{
		Total:         len(linkBaseInfos),
		BaseLinkInfos: linkBaseInfos,
	}, nil
}
