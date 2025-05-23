package logic

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	"shorterurl/link/rpc/internal/svc"
	"shorterurl/link/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StatsGetSingleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStatsGetSingleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StatsGetSingleLogic {
	return &StatsGetSingleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// StatsGetSingle 获取单个短链接监控数据
func (l *StatsGetSingleLogic) StatsGetSingle(in *pb.GetSingleStatsRequest) (*pb.GetSingleStatsResponse, error) {
	// 参数验证
	if in.FullShortUrl == "" {
		return nil, status.Error(codes.InvalidArgument, "短链接不能为空")
	}
	if in.Gid == "" {
		return nil, status.Error(codes.InvalidArgument, "分组标识不能为空")
	}
	if in.StartDate == "" || in.EndDate == "" {
		return nil, status.Error(codes.InvalidArgument, "开始日期和结束日期不能为空")
	}

	// 验证分组是否属于当前用户
	if err := l.checkGroupBelongToUser(in.Gid); err != nil {
		return nil, err
	}

	// 1. 获取基础访问统计数据 (PV、UV、UIP)
	pvUvUip, err := l.svcCtx.RepoManager.LinkAccessStats.FindPvUvUipStatsByShortLink(l.ctx, in.FullShortUrl, in.Gid, in.StartDate, in.EndDate, in.EnableStatus)
	if err != nil {
		l.Logger.Errorf("获取短链接访问量统计失败: %v", err)
		return nil, status.Error(codes.Internal, "获取短链接访问量统计失败")
	}

	// 如果没有任何访问数据，则返回空结果
	if pvUvUip == nil || (pvUvUip.Pv == 0 && pvUvUip.Uv == 0 && pvUvUip.Uip == 0) {
		l.Logger.Info("短链接没有访问数据")
		return &pb.GetSingleStatsResponse{
			Pv:           0,
			Uv:           0,
			Uip:          0,
			Daily:        []*pb.DailyStat{},
			HourStats:    []int32{},
			WeekdayStats: []int32{},
		}, nil
	}

	// 2. 获取每日访问详情
	dailyStats, err := l.svcCtx.RepoManager.LinkAccessStats.ListStatsByShortLink(l.ctx, in.FullShortUrl, in.Gid, in.StartDate, in.EndDate, in.EnableStatus)
	if err != nil {
		l.Logger.Errorf("获取短链接每日访问详情失败: %v", err)
		return nil, status.Error(codes.Internal, "获取短链接每日访问详情失败")
	}

	// 将每日统计数据转换为响应格式
	daily := make([]*pb.DailyStat, 0)
	dateMap := make(map[string]bool)

	// 现有数据放入map
	for _, stat := range dailyStats {
		dateStr := stat.Date.Format("2006-01-02")
		dateMap[dateStr] = true
		daily = append(daily, &pb.DailyStat{
			Date: dateStr,
			Pv:   stat.Pv,
			Uv:   stat.Uv,
			Uip:  stat.Uip,
		})
	}

	// 补充日期范围内没有数据的日期
	startDate, _ := time.Parse("2006-01-02", in.StartDate)
	endDate, _ := time.Parse("2006-01-02", in.EndDate)
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		if _, exists := dateMap[dateStr]; !exists {
			daily = append(daily, &pb.DailyStat{
				Date: dateStr,
				Pv:   0,
				Uv:   0,
				Uip:  0,
			})
		}
	}

	// 3. 获取地区访问详情
	localeStats, err := l.svcCtx.RepoManager.LinkLocaleStats.ListLocaleByShortLink(l.ctx, in.FullShortUrl, in.Gid, in.StartDate, in.EndDate, in.EnableStatus)
	if err != nil {
		l.Logger.Errorf("获取短链接地区访问详情失败: %v", err)
		return nil, status.Error(codes.Internal, "获取短链接地区访问详情失败")
	}

	// 计算地区访问总数
	var localeCnSum int32 = 0
	for _, stat := range localeStats {
		localeCnSum += stat.Cnt
	}

	// 转换地区统计数据
	localeCnStats := make([]*pb.LocaleCnStat, 0, len(localeStats))
	for _, stat := range localeStats {
		// 计算比例并四舍五入到两位小数
		ratio := float64(stat.Cnt) / float64(localeCnSum)
		actualRatio := math.Round(ratio*100.0) / 100.0

		localeCnStats = append(localeCnStats, &pb.LocaleCnStat{
			Locale: stat.Province,
			Cnt:    stat.Cnt,
			Ratio:  actualRatio,
		})
	}

	// 4. 获取小时访问详情
	hourStats, err := l.svcCtx.RepoManager.LinkAccessStats.ListHourStatsByShortLink(l.ctx, in.FullShortUrl, in.Gid, in.StartDate, in.EndDate, in.EnableStatus)
	if err != nil {
		l.Logger.Errorf("获取短链接小时访问详情失败: %v", err)
		return nil, status.Error(codes.Internal, "获取短链接小时访问详情失败")
	}

	// 转换小时统计数据
	hourStatsArray := make([]int32, 24)
	for _, stat := range hourStats {
		if stat.Hour >= 0 && stat.Hour < 24 {
			hourStatsArray[stat.Hour] = stat.Pv
		}
	}

	// 5. 获取高频访问IP详情
	topIpStats, err := l.svcCtx.RepoManager.LinkAccessLogs.ListTopIpByShortLink(l.ctx, in.FullShortUrl, in.Gid, in.StartDate, in.EndDate, in.EnableStatus)
	if err != nil {
		l.Logger.Errorf("获取短链接高频访问IP详情失败: %v", err)
		return nil, status.Error(codes.Internal, "获取短链接高频访问IP详情失败")
	}

	// 计算高频访问IP总数
	ipSum := 0
	for _, stat := range topIpStats {
		var count int
		switch v := stat["count"].(type) {
		case string:
			count, _ = strconv.Atoi(v)
		case int64:
			count = int(v)
		case int:
			count = v
		default:
			count, _ = strconv.Atoi(fmt.Sprintf("%v", v))
		}
		ipSum += count
	}

	// 转换IP统计数据
	topIpStatsResult := make([]*pb.TopIpStat, 0, len(topIpStats))
	for _, stat := range topIpStats {
		var count int
		switch v := stat["count"].(type) {
		case string:
			count, _ = strconv.Atoi(v)
		case int64:
			count = int(v)
		case int:
			count = v
		default:
			count, _ = strconv.Atoi(fmt.Sprintf("%v", v))
		}

		// 计算比例并四舍五入到两位小数
		var actualRatio float64 = 0
		if ipSum > 0 {
			ratio := float64(count) / float64(ipSum)
			actualRatio = math.Round(ratio*100.0) / 100.0
		}

		topIpStatsResult = append(topIpStatsResult, &pb.TopIpStat{
			Ip:    stat["ip"].(string),
			Cnt:   int32(count),
			Ratio: actualRatio,
		})
	}

	// 6. 获取一周访问详情
	weekdayStats, err := l.svcCtx.RepoManager.LinkAccessStats.ListWeekdayStatsByShortLink(l.ctx, in.FullShortUrl, in.Gid, in.StartDate, in.EndDate, in.EnableStatus)
	if err != nil {
		l.Logger.Errorf("获取短链接一周访问详情失败: %v", err)
		return nil, status.Error(codes.Internal, "获取短链接一周访问详情失败")
	}

	// 转换星期统计数据 (1-7 表示周一到周日)
	weekdayStatsArray := make([]int32, 7)
	for _, stat := range weekdayStats {
		if stat.Weekday >= 1 && stat.Weekday <= 7 {
			weekdayStatsArray[stat.Weekday-1] = stat.Pv
		}
	}

	// 7. 获取浏览器访问详情
	browserStats, err := l.svcCtx.RepoManager.LinkBrowserStats.ListBrowserStatsByShortLink(l.ctx, in.FullShortUrl, in.Gid, in.StartDate, in.EndDate, in.EnableStatus)
	if err != nil {
		l.Logger.Errorf("获取短链接浏览器访问详情失败: %v", err)
		return nil, status.Error(codes.Internal, "获取短链接浏览器访问详情失败")
	}

	// 计算浏览器访问总数
	var browserSum int32 = 0
	for _, stat := range browserStats {
		var count int
		switch v := stat["count"].(type) {
		case string:
			count, _ = strconv.Atoi(v)
		case int64:
			count = int(v)
		case int:
			count = v
		default:
			count, _ = strconv.Atoi(fmt.Sprintf("%v", v))
		}
		browserSum += int32(count)
	}

	// 转换浏览器统计数据
	browserStatsResult := make([]*pb.BrowserStat, 0, len(browserStats))
	for _, stat := range browserStats {
		var count int
		switch v := stat["count"].(type) {
		case string:
			count, _ = strconv.Atoi(v)
		case int64:
			count = int(v)
		case int:
			count = v
		default:
			count, _ = strconv.Atoi(fmt.Sprintf("%v", v))
		}

		// 计算比例并四舍五入到两位小数
		var actualRatio float64 = 0
		if browserSum > 0 {
			ratio := float64(count) / float64(browserSum)
			actualRatio = math.Round(ratio*100.0) / 100.0
		}

		browserStatsResult = append(browserStatsResult, &pb.BrowserStat{
			Browser: stat["browser"].(string),
			Cnt:     int32(count),
			Ratio:   actualRatio,
		})
	}

	// 8. 获取操作系统访问详情
	osStats, err := l.svcCtx.RepoManager.LinkOsStats.ListOsStatsByShortLink(l.ctx, in.FullShortUrl, in.Gid, in.StartDate, in.EndDate, in.EnableStatus)
	if err != nil {
		l.Logger.Errorf("获取短链接操作系统访问详情失败: %v", err)
		return nil, status.Error(codes.Internal, "获取短链接操作系统访问详情失败")
	}

	// 计算操作系统访问总数
	var osSum int32 = 0
	for _, stat := range osStats {
		var count int
		switch v := stat["count"].(type) {
		case string:
			count, _ = strconv.Atoi(v)
		case int64:
			count = int(v)
		case int:
			count = v
		default:
			count, _ = strconv.Atoi(fmt.Sprintf("%v", v))
		}
		osSum += int32(count)
	}

	// 转换操作系统统计数据
	osStatsResult := make([]*pb.OSStat, 0, len(osStats))
	for _, stat := range osStats {
		var count int
		switch v := stat["count"].(type) {
		case string:
			count, _ = strconv.Atoi(v)
		case int64:
			count = int(v)
		case int:
			count = v
		default:
			count, _ = strconv.Atoi(fmt.Sprintf("%v", v))
		}

		// 计算比例并四舍五入到两位小数
		var actualRatio float64 = 0
		if osSum > 0 {
			ratio := float64(count) / float64(osSum)
			actualRatio = math.Round(ratio*100.0) / 100.0
		}

		osStatsResult = append(osStatsResult, &pb.OSStat{
			Os:    stat["os"].(string),
			Cnt:   int32(count),
			Ratio: actualRatio,
		})
	}

	// 9. 获取访客类型统计
	uvTypeStats, err := l.svcCtx.RepoManager.LinkAccessLogs.FindUvTypeCntByShortLink(l.ctx, in.FullShortUrl, in.Gid, in.StartDate, in.EndDate, in.EnableStatus)
	if err != nil {
		l.Logger.Errorf("获取短链接访客类型统计失败: %v", err)
		return nil, status.Error(codes.Internal, "获取短链接访客类型统计失败")
	}

	// 转换访客类型统计数据
	oldUserCnt := 0
	newUserCnt := 0
	if len(uvTypeStats) > 0 {
		for _, stat := range uvTypeStats {
			uvType, ok := stat["uv_type"].(string)
			if !ok {
				continue
			}

			var count int
			switch v := stat["count"].(type) {
			case string:
				count, _ = strconv.Atoi(v)
			case int64:
				count = int(v)
			case int:
				count = v
			default:
				count, _ = strconv.Atoi(fmt.Sprintf("%v", v))
			}

			if uvType == "新访客" {
				newUserCnt = count
			} else if uvType == "旧访客" {
				oldUserCnt = count
			}
		}
	}

	// 计算访客类型比例
	uvSum := oldUserCnt + newUserCnt
	var oldRatio, newRatio float64 = 0, 0
	if uvSum > 0 {
		oldRatio = float64(oldUserCnt) / float64(uvSum)
		newRatio = float64(newUserCnt) / float64(uvSum)
	}
	actualOldRatio := math.Round(oldRatio*100.0) / 100.0
	actualNewRatio := math.Round(newRatio*100.0) / 100.0

	uvTypeStatsResult := []*pb.UvTypeStat{
		{
			UvType: "newUser",
			Cnt:    int32(newUserCnt),
			Ratio:  actualNewRatio,
		},
		{
			UvType: "oldUser",
			Cnt:    int32(oldUserCnt),
			Ratio:  actualOldRatio,
		},
	}

	// 10. 获取设备类型访问详情
	deviceStats, err := l.svcCtx.RepoManager.LinkDeviceStats.ListDeviceStatsByShortLink(l.ctx, in.FullShortUrl, in.Gid, in.StartDate, in.EndDate, in.EnableStatus)
	if err != nil {
		l.Logger.Errorf("获取短链接设备类型访问详情失败: %v", err)
		return nil, status.Error(codes.Internal, "获取短链接设备类型访问详情失败")
	}

	// 计算设备类型访问总数
	var deviceSum int32 = 0
	for _, stat := range deviceStats {
		deviceSum += stat.Cnt
	}

	// 转换设备类型统计数据
	deviceStatsResult := make([]*pb.DeviceStat, 0, len(deviceStats))
	for _, stat := range deviceStats {
		// 计算比例并四舍五入到两位小数
		var actualRatio float64 = 0
		if deviceSum > 0 {
			ratio := float64(stat.Cnt) / float64(deviceSum)
			actualRatio = math.Round(ratio*100.0) / 100.0
		}

		deviceStatsResult = append(deviceStatsResult, &pb.DeviceStat{
			Device: stat.Device,
			Cnt:    stat.Cnt,
			Ratio:  actualRatio,
		})
	}

	// 11. 获取网络类型访问详情
	networkStats, err := l.svcCtx.RepoManager.LinkNetworkStats.ListNetworkStatsByShortLink(l.ctx, in.FullShortUrl, in.Gid, in.StartDate, in.EndDate, in.EnableStatus)
	if err != nil {
		l.Logger.Errorf("获取短链接网络类型访问详情失败: %v", err)
		return nil, status.Error(codes.Internal, "获取短链接网络类型访问详情失败")
	}

	// 计算网络类型访问总数
	var networkSum int32 = 0
	for _, stat := range networkStats {
		networkSum += stat.Cnt
	}

	// 转换网络类型统计数据
	networkStatsResult := make([]*pb.NetworkStat, 0, len(networkStats))
	for _, stat := range networkStats {
		// 计算比例并四舍五入到两位小数
		var actualRatio float64 = 0
		if networkSum > 0 {
			ratio := float64(stat.Cnt) / float64(networkSum)
			actualRatio = math.Round(ratio*100.0) / 100.0
		}

		networkStatsResult = append(networkStatsResult, &pb.NetworkStat{
			Network: stat.Network,
			Cnt:     stat.Cnt,
			Ratio:   actualRatio,
		})
	}

	// 构建并返回结果
	return &pb.GetSingleStatsResponse{
		Pv:            pvUvUip.Pv,
		Uv:            pvUvUip.Uv,
		Uip:           pvUvUip.Uip,
		Daily:         daily,
		LocaleCnStats: localeCnStats,
		HourStats:     hourStatsArray,
		TopIpStats:    topIpStatsResult,
		WeekdayStats:  weekdayStatsArray,
		BrowserStats:  browserStatsResult,
		OsStats:       osStatsResult,
		UvTypeStats:   uvTypeStatsResult,
		DeviceStats:   deviceStatsResult,
		NetworkStats:  networkStatsResult,
	}, nil
}

// checkGroupBelongToUser 检查分组是否属于当前用户
func (l *StatsGetSingleLogic) checkGroupBelongToUser(gid string) error {
	// 获取当前登录用户
	username, err := l.svcCtx.RepoManager.GetCurrentUsername(l.ctx)
	if err != nil {
		return status.Error(codes.Unauthenticated, "用户未登录")
	}

	// 检查分组是否属于该用户
	exist, err := l.svcCtx.RepoManager.Group.CheckGroupBelongToUser(l.ctx, gid, username)
	if err != nil {
		l.Logger.Errorf("检查分组归属失败: %v", err)
		return status.Error(codes.Internal, "检查分组归属失败")
	}

	if !exist {
		return status.Error(codes.PermissionDenied, "用户信息与分组标识不匹配")
	}

	return nil
}
