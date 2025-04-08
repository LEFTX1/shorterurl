package stats

import (
	"context"
	"shorterurl/link/rpc/shortlinkservice"
	"shorterurl/user/api/internal/svc"
	"shorterurl/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"
)

type ShortLinkStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取单个短链接监控数据
func NewShortLinkStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShortLinkStatsLogic {
	return &ShortLinkStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShortLinkStatsLogic) ShortLinkStats(req *types.ShortLinkStatsReq) (resp *types.ShortLinkStatsRespDTO, err error) {
	// 获取当前用户信息
	userInfo := l.ctx.Value(types.UserContextKey).(*types.UserInfo)

	// 创建新的上下文并添加用户信息
	ctx := metadata.NewOutgoingContext(l.ctx, metadata.Pairs(
		"username", userInfo.Username,
	))

	// 调用link RPC服务获取短链接统计数据
	result, err := l.svcCtx.LinkRpc.StatsGetSingle(ctx, &shortlinkservice.GetSingleStatsRequest{
		FullShortUrl: req.FullShortUrl,
		Gid:          req.Gid,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		EnableStatus: int32(req.EnableStatus),
	})
	if err != nil {
		logx.Errorf("获取短链接统计数据失败: %v", err)
		return nil, err
	}

	// 构建响应
	resp = &types.ShortLinkStatsRespDTO{}

	// 转换PV/UV/UIP统计数据
	pvUvUipStats := make([]types.PvUvUipStats, 0)
	for _, stat := range result.Daily {
		pvUvUipStats = append(pvUvUipStats, types.PvUvUipStats{
			Date: stat.Date,
			Pv:   int64(stat.Pv),
			Uv:   int64(stat.Uv),
			Uip:  int64(stat.Uip),
		})
	}

	resp.PvUvUipStatsList = pvUvUipStats

	// 设置总体统计数据
	resp.OverallPvUvUipStats = types.PvUvUipStats{
		Date: "总计",
		Pv:   int64(result.Pv),
		Uv:   int64(result.Uv),
		Uip:  int64(result.Uip),
	}

	// 转换地域统计数据
	localeCnStats := make([]types.LocaleCnStat, 0)
	for _, stat := range result.LocaleCnStats {
		localeCnStats = append(localeCnStats, types.LocaleCnStat{
			Locale: stat.Locale,
			Cnt:    int64(stat.Cnt),
			Ratio:  stat.Ratio,
		})
	}
	resp.LocaleCnStats = localeCnStats

	// 设置小时访问详情
	hourStats := make([]int, len(result.HourStats))
	for i, value := range result.HourStats {
		hourStats[i] = int(value)
	}
	resp.HourStats = hourStats

	// 转换高频访问IP统计
	topIpStats := make([]types.TopIpStat, 0)
	for _, stat := range result.TopIpStats {
		topIpStats = append(topIpStats, types.TopIpStat{
			Ip:    stat.Ip,
			Cnt:   int64(stat.Cnt),
			Ratio: stat.Ratio,
		})
	}
	resp.TopIpStats = topIpStats

	// 设置一周访问详情
	weekdayStats := make([]int, len(result.WeekdayStats))
	for i, value := range result.WeekdayStats {
		weekdayStats[i] = int(value)
	}
	resp.WeekdayStats = weekdayStats

	// 转换浏览器统计
	browserStats := make([]types.BrowserStat, 0)
	for _, stat := range result.BrowserStats {
		browserStats = append(browserStats, types.BrowserStat{
			Browser: stat.Browser,
			Cnt:     int64(stat.Cnt),
			Ratio:   stat.Ratio,
		})
	}
	resp.BrowserStats = browserStats

	// 转换操作系统统计
	osStats := make([]types.OsStat, 0)
	for _, stat := range result.OsStats {
		osStats = append(osStats, types.OsStat{
			Os:    stat.Os,
			Cnt:   int64(stat.Cnt),
			Ratio: stat.Ratio,
		})
	}
	resp.OsStats = osStats

	// 转换访客类型统计
	uvTypeStats := make([]types.UvTypeStat, 0)
	for _, stat := range result.UvTypeStats {
		uvTypeStats = append(uvTypeStats, types.UvTypeStat{
			UvType: stat.UvType,
			Cnt:    int64(stat.Cnt),
			Ratio:  stat.Ratio,
		})
	}
	resp.UvTypeStats = uvTypeStats

	// 转换设备统计
	deviceStats := make([]types.DeviceStat, 0)
	for _, stat := range result.DeviceStats {
		deviceStats = append(deviceStats, types.DeviceStat{
			Device: stat.Device,
			Cnt:    int64(stat.Cnt),
			Ratio:  stat.Ratio,
		})
	}
	resp.DeviceStats = deviceStats

	// 转换网络统计
	networkStats := make([]types.NetworkStat, 0)
	for _, stat := range result.NetworkStats {
		networkStats = append(networkStats, types.NetworkStat{
			Network: stat.Network,
			Cnt:     int64(stat.Cnt),
			Ratio:   stat.Ratio,
		})
	}
	resp.NetworkStats = networkStats

	return resp, nil
}
