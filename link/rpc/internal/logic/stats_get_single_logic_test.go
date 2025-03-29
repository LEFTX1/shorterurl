package logic_test

import (
	"context"
	"flag"
	"fmt"
	"shorterurl/link/rpc/internal/config"
	"shorterurl/link/rpc/internal/logic"
	"shorterurl/link/rpc/internal/model"
	"shorterurl/link/rpc/internal/svc"
	"shorterurl/link/rpc/pb"
	"sync"
	"testing"
	"time"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	testOnce sync.Once
	testCtx  context.Context
	testSvc  *svc.ServiceContext
)

func initTest(t *testing.T) (*svc.ServiceContext, context.Context) {
	testOnce.Do(func() {
		// 检查标志是否已定义
		var configFile string
		if f := flag.Lookup("f"); f != nil {
			configFile = f.Value.String()
		} else {
			// 如果未定义，则定义它
			flag.StringVar(&configFile, "f", "../../etc/link.yaml", "配置文件路径")
		}
		// 确保 flag.Parse() 只被调用一次
		if !flag.Parsed() {
			flag.Parse()
		}

		var c config.Config
		// 使用获取到的 configFile 变量
		conf.MustLoad(configFile, &c)

		testSvc = svc.NewServiceContext(c)
		logx.Disable()

		// 创建测试用户上下文
		testCtx = context.WithValue(context.Background(), "username", "test_user")
	})

	return testSvc, testCtx
}

// 为测试准备短链接相关数据
func prepareTestData(t *testing.T, svcCtx *svc.ServiceContext, ctx context.Context) (string, string, func()) {
	// 生成随机的短链接后缀以避免冲突，确保不超过8个字符
	randomID := fmt.Sprintf("t%d", time.Now().UnixNano()%10000) // 限制为最多5个字符 (t + 最多4位数字)
	fullShortUrl := fmt.Sprintf("test.example.com/%s", randomID)
	gid := "test-group-stats"

	// 1. 创建测试用户
	testUser := &model.User{
		Username:   "test_user",
		Password:   "password",
		RealName:   "测试用户",
		Phone:      "13800138000",
		Mail:       "test@example.com",
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		DelFlag:    0,
	}

	err := svcCtx.RepoManager.User.Create(ctx, testUser)
	if err != nil {
		// 如果用户已存在，忽略错误
		t.Logf("创建测试用户失败(可能已存在): %v", err)
	}

	// 2. 创建测试分组
	testGroup := &model.Group{
		Gid:        gid,
		Name:       "测试分组Stats",
		Username:   "test_user", // 分片键
		SortOrder:  1,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		DelFlag:    0,
	}

	err = svcCtx.RepoManager.Group.Create(ctx, testGroup)
	if err != nil {
		// 如果分组已存在，忽略错误
		t.Logf("创建测试分组失败(可能已存在): %v", err)
	}

	// 3. 创建测试短链接
	testLink := &model.Link{
		Domain:        "test.example.com",
		ShortUri:      randomID,
		FullShortUrl:  fullShortUrl,
		OriginUrl:     "https://github.com/zeromicro/go-zero",
		Gid:           gid, // 分片键
		Favicon:       "",
		EnableStatus:  0,
		CreatedType:   0,
		ValidDateType: 0,
		ValidDate:     time.Now().AddDate(1, 0, 0),
		Describe:      "测试短链接",
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
		DelFlag:       0,
	}

	err = svcCtx.RepoManager.Link.Create(ctx, testLink)
	if err != nil {
		t.Fatalf("创建测试短链接失败: %v", err)
	}

	// 4. 创建短链接跳转记录
	testLinkGoto := &model.LinkGoto{
		Gid:          gid,
		FullShortUrl: fullShortUrl, // 分片键
	}

	err = svcCtx.RepoManager.LinkGoto.Create(ctx, testLinkGoto)
	if err != nil {
		t.Fatalf("创建测试短链接跳转记录失败: %v", err)
	}

	// 5. 创建访问统计数据
	// 创建三天的访问数据
	now := time.Now()
	dates := []time.Time{
		now.AddDate(0, 0, -2), // 前天
		now.AddDate(0, 0, -1), // 昨天
		now,                   // 今天
	}

	for _, date := range dates {
		// 访问统计
		accessStat := &model.LinkAccessStats{
			FullShortUrl: fullShortUrl,
			Date:         date,
			PV:           10,
			UV:           5,
			UIP:          3,
			Hour:         date.Hour(),
			Weekday:      int(date.Weekday()),
			CreateTime:   time.Now(),
			UpdateTime:   time.Now(),
			DelFlag:      0,
		}

		err = svcCtx.RepoManager.GetCommonDB().Create(accessStat).Error
		if err != nil {
			t.Fatalf("创建访问统计数据失败: %v", err)
		}

		// 地区统计
		localeStat := &model.LinkLocaleStats{
			FullShortUrl: fullShortUrl,
			Date:         date,
			Cnt:          5,
			Province:     "北京",
			City:         "北京",
			Adcode:       "110000",
			Country:      "CN",
			CreateTime:   time.Now(),
			UpdateTime:   time.Now(),
			DelFlag:      0,
		}

		err = svcCtx.RepoManager.GetCommonDB().Create(localeStat).Error
		if err != nil {
			t.Fatalf("创建地区统计数据失败: %v", err)
		}

		// 浏览器统计
		browserStat := &model.LinkBrowserStats{
			FullShortUrl: fullShortUrl,
			Date:         date,
			Cnt:          5,
			Browser:      "Chrome",
			CreateTime:   time.Now(),
			UpdateTime:   time.Now(),
			DelFlag:      0,
		}

		err = svcCtx.RepoManager.GetCommonDB().Create(browserStat).Error
		if err != nil {
			t.Fatalf("创建浏览器统计数据失败: %v", err)
		}

		// 操作系统统计
		osStat := &model.LinkOsStats{
			FullShortUrl: fullShortUrl,
			Date:         date,
			Cnt:          5,
			Os:           "Windows",
			CreateTime:   time.Now(),
			UpdateTime:   time.Now(),
			DelFlag:      0,
		}

		err = svcCtx.RepoManager.GetCommonDB().Create(osStat).Error
		if err != nil {
			t.Fatalf("创建操作系统统计数据失败: %v", err)
		}

		// 设备统计
		deviceStat := &model.LinkDeviceStats{
			FullShortUrl: fullShortUrl,
			Date:         date,
			Cnt:          5,
			Device:       "PC",
			CreateTime:   time.Now(),
			UpdateTime:   time.Now(),
			DelFlag:      0,
		}

		err = svcCtx.RepoManager.GetCommonDB().Create(deviceStat).Error
		if err != nil {
			t.Fatalf("创建设备统计数据失败: %v", err)
		}

		// 网络统计
		networkStat := &model.LinkNetworkStats{
			FullShortUrl: fullShortUrl,
			Date:         date,
			Cnt:          5,
			Network:      "WIFI",
			CreateTime:   time.Now(),
			UpdateTime:   time.Now(),
			DelFlag:      0,
		}

		err = svcCtx.RepoManager.GetCommonDB().Create(networkStat).Error
		if err != nil {
			t.Fatalf("创建网络统计数据失败: %v", err)
		}
	}

	// 6. 创建访问日志数据，用于高频IP统计和访客类型
	for i := 0; i < 5; i++ {
		accessLog := &model.LinkAccessLog{
			FullShortUrl: fullShortUrl,
			User:         fmt.Sprintf("user%d", i%2),         // 两个不同的用户，模拟新旧访客
			IP:           fmt.Sprintf("192.168.1.%d", i%3+1), // 3个不同IP
			Browser:      "Chrome",
			Os:           "Windows",
			Network:      "WIFI",
			Device:       "PC",
			Locale:       "北京",
			CreateTime:   time.Now(),
			UpdateTime:   time.Now(),
			DelFlag:      0,
		}

		err = svcCtx.RepoManager.GetCommonDB().Create(accessLog).Error
		if err != nil {
			t.Fatalf("创建访问日志数据失败: %v", err)
		}
	}

	// 返回清理函数
	return fullShortUrl, gid, func() {
		// 清理所有测试数据
		// 1. 清理访问日志
		svcCtx.RepoManager.GetCommonDB().Where("full_short_url = ?", fullShortUrl).Delete(&model.LinkAccessLog{})

		// 2. 清理各种统计数据
		svcCtx.RepoManager.GetCommonDB().Where("full_short_url = ?", fullShortUrl).Delete(&model.LinkAccessStats{})
		svcCtx.RepoManager.GetCommonDB().Where("full_short_url = ?", fullShortUrl).Delete(&model.LinkLocaleStats{})
		svcCtx.RepoManager.GetCommonDB().Where("full_short_url = ?", fullShortUrl).Delete(&model.LinkBrowserStats{})
		svcCtx.RepoManager.GetCommonDB().Where("full_short_url = ?", fullShortUrl).Delete(&model.LinkOsStats{})
		svcCtx.RepoManager.GetCommonDB().Where("full_short_url = ?", fullShortUrl).Delete(&model.LinkDeviceStats{})
		svcCtx.RepoManager.GetCommonDB().Where("full_short_url = ?", fullShortUrl).Delete(&model.LinkNetworkStats{})

		// 3. 清理短链接跳转记录
		svcCtx.RepoManager.LinkGoto.Delete(ctx, fullShortUrl)

		// 4. 清理短链接 - 获取ID后删除
		var link model.Link
		err := svcCtx.RepoManager.GetLinkDB().Where("full_short_url = ? AND gid = ?", fullShortUrl, gid).First(&link).Error
		if err == nil && link.ID > 0 {
			svcCtx.RepoManager.Link.Delete(ctx, link.ID, gid)
		}

		// 注意：不删除测试用户和分组，以便复用
	}
}

// TestStatsGetSingle_Success 测试正常获取单个短链接统计数据
func TestStatsGetSingle_Success(t *testing.T) {
	svcCtx, ctx := initTest(t)

	// 准备测试数据
	fullShortUrl, gid, cleanup := prepareTestData(t, svcCtx, ctx)
	defer cleanup()

	// 测试获取统计数据
	l := logic.NewStatsGetSingleLogic(ctx, svcCtx)
	now := time.Now()
	startDate := now.AddDate(0, 0, -3).Format("2006-01-02") // 3天前
	endDate := now.Format("2006-01-02")                     // 今天

	resp, err := l.StatsGetSingle(&pb.GetSingleStatsRequest{
		FullShortUrl: fullShortUrl,
		Gid:          gid,
		StartDate:    startDate,
		EndDate:      endDate,
		EnableStatus: 0,
	})

	if err != nil {
		t.Fatalf("获取单个短链接统计数据失败: %v", err)
	}

	if resp == nil {
		t.Fatal("响应为空")
	}

	// 验证基本数据
	if resp.Pv <= 0 {
		t.Errorf("期望PV > 0，实际: %d", resp.Pv)
	}

	if resp.Uv <= 0 {
		t.Errorf("期望UV > 0，实际: %d", resp.Uv)
	}

	if resp.Uip <= 0 {
		t.Errorf("期望UIP > 0，实际: %d", resp.Uip)
	}

	// 验证日常统计数据
	if len(resp.Daily) == 0 {
		t.Error("每日统计数据为空")
	}

	// 验证地区统计数据
	if len(resp.LocaleCnStats) == 0 {
		t.Error("地区统计数据为空")
	}

	// 验证小时统计数据
	hourSum := int32(0)
	for _, h := range resp.HourStats {
		hourSum += h
	}
	if hourSum <= 0 {
		t.Error("小时统计数据总和应大于0")
	}

	// 验证星期统计数据
	weekdaySum := int32(0)
	for _, w := range resp.WeekdayStats {
		weekdaySum += w
	}
	if weekdaySum <= 0 {
		t.Error("星期统计数据总和应大于0")
	}

	// 验证浏览器统计数据
	if len(resp.BrowserStats) == 0 {
		t.Error("浏览器统计数据为空")
	}

	// 验证操作系统统计数据
	if len(resp.OsStats) == 0 {
		t.Error("操作系统统计数据为空")
	}

	// 验证访客类型统计数据
	if len(resp.UvTypeStats) == 0 {
		t.Error("访客类型统计数据为空")
	}

	// 验证设备统计数据
	if len(resp.DeviceStats) == 0 {
		t.Error("设备统计数据为空")
	}

	// 验证网络统计数据
	if len(resp.NetworkStats) == 0 {
		t.Error("网络统计数据为空")
	}

	t.Logf("成功获取单个短链接统计数据: PV=%d, UV=%d, UIP=%d", resp.Pv, resp.Uv, resp.Uip)
}

// TestStatsGetSingle_InvalidParams 测试无效参数
func TestStatsGetSingle_InvalidParams(t *testing.T) {
	svcCtx, ctx := initTest(t)

	l := logic.NewStatsGetSingleLogic(ctx, svcCtx)

	// 测试空短链接
	_, err := l.StatsGetSingle(&pb.GetSingleStatsRequest{
		FullShortUrl: "",
		Gid:          "test-group",
		StartDate:    "2023-01-01",
		EndDate:      "2023-01-31",
		EnableStatus: 0,
	})

	if err == nil {
		t.Error("期望空短链接参数校验失败，但实际成功")
	}

	// 测试空分组ID
	_, err = l.StatsGetSingle(&pb.GetSingleStatsRequest{
		FullShortUrl: "test.example.com/test",
		Gid:          "",
		StartDate:    "2023-01-01",
		EndDate:      "2023-01-31",
		EnableStatus: 0,
	})

	if err == nil {
		t.Error("期望空分组ID参数校验失败，但实际成功")
	}

	// 测试空日期范围
	_, err = l.StatsGetSingle(&pb.GetSingleStatsRequest{
		FullShortUrl: "test.example.com/test",
		Gid:          "test-group",
		StartDate:    "",
		EndDate:      "",
		EnableStatus: 0,
	})

	if err == nil {
		t.Error("期望空日期范围参数校验失败，但实际成功")
	}
}

// TestStatsGetSingle_NoData 测试没有数据的情况
func TestStatsGetSingle_NoData(t *testing.T) {
	svcCtx, ctx := initTest(t)

	// 生成一个不存在的短链接
	nonExistingUrl := fmt.Sprintf("test.example.com/nonexist%d", time.Now().UnixNano())

	l := logic.NewStatsGetSingleLogic(ctx, svcCtx)
	resp, err := l.StatsGetSingle(&pb.GetSingleStatsRequest{
		FullShortUrl: nonExistingUrl,
		Gid:          "test-group-stats", // 使用已存在的分组ID
		StartDate:    "2023-01-01",
		EndDate:      "2023-01-31",
		EnableStatus: 0,
	})

	if err != nil {
		t.Fatalf("获取不存在链接的统计数据应该返回空结果而不是错误: %v", err)
	}

	if resp == nil {
		t.Fatal("响应不应为空，应返回空数据")
	}

	// 验证返回的是空数据
	if resp.Pv != 0 || resp.Uv != 0 || resp.Uip != 0 {
		t.Errorf("期望PV=UV=UIP=0，实际: PV=%d, UV=%d, UIP=%d", resp.Pv, resp.Uv, resp.Uip)
	}

	if len(resp.Daily) != 0 {
		t.Errorf("期望每日统计为空数组，实际长度: %d", len(resp.Daily))
	}
}

// TestStatsGetSingle_PermissionDenied 测试权限不足的情况
func TestStatsGetSingle_PermissionDenied(t *testing.T) {
	svcCtx, ctx := initTest(t)

	// 创建一个无权限上下文（没有用户名）
	noAuthCtx := context.Background()

	l := logic.NewStatsGetSingleLogic(noAuthCtx, svcCtx)
	_, err := l.StatsGetSingle(&pb.GetSingleStatsRequest{
		FullShortUrl: "test.example.com/test",
		Gid:          "test-group-stats",
		StartDate:    "2023-01-01",
		EndDate:      "2023-01-31",
		EnableStatus: 0,
	})

	if err == nil {
		t.Error("期望未登录用户访问时返回错误，但实际成功")
	}

	// 使用一个不存在的分组ID（但用户存在）
	l = logic.NewStatsGetSingleLogic(ctx, svcCtx)
	_, err = l.StatsGetSingle(&pb.GetSingleStatsRequest{
		FullShortUrl: "test.example.com/test",
		Gid:          "non-existing-group",
		StartDate:    "2023-01-01",
		EndDate:      "2023-01-31",
		EnableStatus: 0,
	})

	if err == nil {
		t.Error("期望不存在的分组ID返回错误，但实际成功")
	}
}
