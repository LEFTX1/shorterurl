package recycle

import (
	"context"
	"net/url"
	"shorterurl/link/rpc/shortlinkservice"
	"shorterurl/user/api/internal/svc"
	"shorterurl/user/api/internal/types"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecycleBinPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRecycleBinPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecycleBinPageLogic {
	return &RecycleBinPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RecycleBinPageLogic) RecycleBinPage(req *types.RecycleBinPageReq) (resp *types.RecycleBinPageResp, err error) {
	// 获取请求中的分组参数
	gid := req.Gid
	l.Logger.Infof("查询回收站短链接，分组标识: %s", gid)

	// 调用link RPC服务获取回收站短链接列表
	result, err := l.svcCtx.LinkRpc.RecycleBinPage(l.ctx, &shortlinkservice.PageRecycleBinShortLinkRequest{
		Gid:     gid,
		Current: int32(req.Current),
		Size:    int32(req.Size),
	})
	if err != nil {
		logx.Errorf("获取回收站短链接列表失败: %v", err)
		return nil, err
	}

	// 构建响应
	resp = &types.RecycleBinPageResp{
		Total:   int64(result.Total),
		Size:    int(result.Size),
		Current: int(result.Current),
		Records: make([]types.ShortLinkPageRecordDTO, 0, len(result.Records)),
	}

	// 转换短链接记录数据
	for _, record := range result.Records {
		// 提取shortUri
		shortUri := extractShortUri(record.FullShortUrl)

		// 确定有效期类型
		validDateType := 0
		if record.ValidDate != "" {
			validDateType = 1
		}

		// 构建短链接记录
		item := types.ShortLinkPageRecordDTO{
			Domain:        record.Domain,
			FullShortUrl:  record.FullShortUrl,
			ShortUri:      shortUri,
			OriginUrl:     record.OriginUrl,
			Gid:           record.Gid,
			CreateTime:    record.CreateTime,
			Describe:      record.Describe,
			ValidDate:     record.ValidDate,
			ValidDateType: validDateType,
			TotalPv:       int64(record.TotalPv),
			TotalUv:       int64(record.TotalUv),
			TotalUip:      int64(record.TotalUip),
			// 设置默认值
			Id:           0,
			Favicon:      "https://cdn-icons-png.flaticon.com/512/8763/8763935.png", // 默认图标
			EnableStatus: 1,                                                         // 回收站中的链接EnableStatus=1
			TodayPv:      0,                                                         // 今日数据默认为0
			TodayUv:      0,
			TodayUip:     0,
		}

		resp.Records = append(resp.Records, item)
	}

	return resp, nil
}

// 从完整短链接中提取短链接URI部分
func extractShortUri(fullShortUrl string) string {
	if fullShortUrl == "" {
		return ""
	}

	// 使用标准库解析URL
	parsedUrl, err := url.Parse(fullShortUrl)
	if err != nil {
		// 如果解析失败，使用基础字符串处理
		parts := strings.Split(fullShortUrl, "/")
		if len(parts) > 0 {
			return parts[len(parts)-1]
		}
		return ""
	}

	// 获取路径的最后一部分
	path := parsedUrl.Path
	if path == "" {
		return ""
	}

	// 去除开头的斜杠
	path = strings.TrimPrefix(path, "/")

	// 如果路径中还有斜杠，只取最后一部分
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}
