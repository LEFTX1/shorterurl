package stats

import (
	"context"
	"shorterurl/link/rpc/shortlinkservice"
	"shorterurl/user/api/internal/svc"
	"shorterurl/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"
)

type ShortLinkAccessRecordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 短链接访问记录查询
func NewShortLinkAccessRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShortLinkAccessRecordLogic {
	return &ShortLinkAccessRecordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShortLinkAccessRecordLogic) ShortLinkAccessRecord(req *types.ShortLinkAccessRecordReq) (resp *types.AccessRecordPageResp, err error) {
	// 获取当前用户信息
	userInfo := l.ctx.Value(types.UserContextKey).(*types.UserInfo)

	// 创建新的上下文并添加用户信息
	ctx := metadata.NewOutgoingContext(l.ctx, metadata.Pairs(
		"username", userInfo.Username,
	))

	// 调用link RPC服务获取短链接访问记录
	result, err := l.svcCtx.LinkRpc.StatsAccessRecordQuery(ctx, &shortlinkservice.AccessRecordQueryRequest{
		FullShortUrl: req.FullShortUrl,
		Gid:          req.Gid,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		EnableStatus: int32(req.EnableStatus),
		Current:      req.Current,
		Size:         req.Size,
	})
	if err != nil {
		logx.Errorf("获取短链接访问记录失败: %v", err)
		return nil, err
	}

	// 构建响应
	resp = &types.AccessRecordPageResp{
		Total:   result.Total,
		Size:    result.Size,
		Current: result.Current,
		Records: make([]types.AccessRecord, 0, len(result.Records)),
	}

	// 转换访问记录数据
	for _, record := range result.Records {
		// 将RPC结果映射到API响应
		resp.Records = append(resp.Records, types.AccessRecord{
			Ip:         record.Ip,
			Browser:    record.Browser,
			Os:         record.Os,
			Network:    record.Network,
			Device:     record.Device,
			Locale:     record.Locale,
			AccessTime: record.CreateTime,
			// API类型中没有UvType和User字段，所以这些数据在这里会丢失
		})
	}

	return resp, nil
}
