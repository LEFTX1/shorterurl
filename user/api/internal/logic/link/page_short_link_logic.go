package link

import (
	"context"

	"shorterurl/user/api/internal/svc"
	"shorterurl/user/api/internal/types"

	"errors"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"

	"shorterurl/link/rpc/shortlinkservice"
)

type PageShortLinkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 分页查询短链接
func NewPageShortLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageShortLinkLogic {
	return &PageShortLinkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PageShortLinkLogic) PageShortLink(req *types.PageLinkReq) (resp *types.PageLinkResp, err error) {
	// 获取用户信息
	userInfo, ok := types.GetUserFromCtx(l.ctx)
	if !ok {
		return nil, errors.New("未找到用户信息")
	}

	// 调用 RPC 服务
	rpcReq := &shortlinkservice.PageShortLinkRequest{
		Gid:     req.Gid,
		Current: int32(req.Current),
		Size:    int32(req.Size),
	}

	// 添加元数据
	ctx := metadata.AppendToOutgoingContext(l.ctx, "username", userInfo.Username)
	rpcResp, err := l.svcCtx.LinkRpc.ShortLinkPage(ctx, rpcReq)
	if err != nil {
		l.Logger.Errorf("分页查询短链接失败 username: %s, gid: %s, error: %v", userInfo.Username, req.Gid, err)
		return nil, err
	}

	// 构建响应
	records := make([]types.ShortLinkRecord, 0, len(rpcResp.Records))
	for _, record := range rpcResp.Records {
		records = append(records, types.ShortLinkRecord{
			FullShortUrl: record.FullShortUrl,
			OriginUrl:    record.OriginUrl,
			Domain:       record.Domain,
			Gid:          record.Gid,
			CreateTime:   record.CreateTime,
			ValidDate:    record.ValidDate,
			Describe:     record.Describe,
			TotalPv:      int64(record.TotalPv),
			TotalUv:      int64(record.TotalUv),
			TotalUip:     int64(record.TotalUip),
			// 其他统计字段暂时不需要填充
		})
	}

	return &types.PageLinkResp{
		Records: records,
		Total:   int64(rpcResp.Total),
		Size:    int(rpcResp.Size),
		Current: int(rpcResp.Current),
	}, nil
}
