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
	groupTestOnce sync.Once
	groupTestCtx  context.Context
	groupTestSvc  *svc.ServiceContext
)

// initGroupTest 设置测试环境
func initGroupTest(t *testing.T) (*svc.ServiceContext, context.Context) {
	groupTestOnce.Do(func() {
		// 检查标志是否已定义
		var configFile string
		if f := flag.Lookup("f"); f != nil {
			configFile = f.Value.String()
		} else {
			// 如果未定义，则定义它
			flag.StringVar(&configFile, "f", "../../etc/link.yaml", "配置文件路径")
		}
		// 确保 flag.Parse() 在所有标志定义之后只调用一次，
		// 通常在 TestMain 中调用是最佳实践，这里暂时保留在 Once 中，
		// 但要注意如果其他测试文件也调用 flag.Parse() 可能引发问题。
		// 为了安全起见，可以在 TestMain 中统一处理 flag.Parse()。
		// 如果没有 TestMain，则需要确保 Parse 只执行一次。
		// 这里假设 TestMain 中没有调用 flag.Parse()。
		if !flag.Parsed() {
			flag.Parse()
		}

		var c config.Config
		// 使用获取到的 configFile 变量
		conf.MustLoad(configFile, &c)

		groupTestSvc = svc.NewServiceContext(c)
		logx.Disable()

		// 创建测试用户上下文
		groupTestCtx = context.WithValue(context.Background(), "username", "test_user")

	})

	return groupTestSvc, groupTestCtx
}

// 为测试准备分组相关数据
func prepareGroupTestData(t *testing.T, svcCtx *svc.ServiceContext, ctx context.Context) (string, func()) {
	// 生成不超过32个字符的唯一分组ID
	gid := fmt.Sprintf("test-group-%d", time.Now().UnixNano()%100000)

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
		t.Fatalf("创建测试分组失败: %v", err)
	}

	// 3. 创建多个测试短链接
	for i := 0; i < 3; i++ {
		// 生成不超过8个字符的短链接URI
		shortUri := fmt.Sprintf("g%d", i%1000)
		fullShortUrl := fmt.Sprintf("test.example.com/%s", shortUri)

		testLink := &model.Link{
			Domain:        "test.example.com",
			ShortUri:      shortUri,
			FullShortUrl:  fullShortUrl,
			OriginUrl:     fmt.Sprintf("https://github.com/zeromicro/go-zero/page%d", i),
			Gid:           gid, // 分片键
			Favicon:       "",
			EnableStatus:  0,
			CreatedType:   0,
			ValidDateType: 0,
			ValidDate:     time.Now().AddDate(1, 0, 0),
			Describe:      fmt.Sprintf("测试短链接-%d", i),
			CreateTime:    time.Now(),
			UpdateTime:    time.Now(),
			DelFlag:       0,
		}

		err = svcCtx.RepoManager.Link.Create(ctx, testLink)
		if err != nil {
			t.Fatalf("创建测试短链接失败: %v", err)
		}

		// 创建短链接跳转记录
		testLinkGoto := &model.LinkGoto{
			Gid:          gid,
			FullShortUrl: fullShortUrl, // 分片键
		}

		err = svcCtx.RepoManager.LinkGoto.Create(ctx, testLinkGoto)
		if err != nil {
			t.Fatalf("创建测试短链接跳转记录失败: %v", err)
		}

		// 为每个短链接创建访问统计数据
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

		// 创建访问日志数据，用于高频IP统计和访客类型
		for j := 0; j < 5; j++ {
			accessLog := &model.LinkAccessLog{
				FullShortUrl: fullShortUrl,
				User:         fmt.Sprintf("user%d", j%2),         // 两个不同的用户，模拟新旧访客
				IP:           fmt.Sprintf("192.168.1.%d", j%3+1), // 3个不同IP
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
	}

	// 返回清理函数
	return gid, func() {
		// 查找并删除该分组下的所有短链接
		links, err := svcCtx.RepoManager.Link.FindByGidWithCondition(ctx, gid, map[string]interface{}{}, 1, 100)
		if err != nil {
			t.Logf("查询分组下短链接失败: %v", err)
		} else {
			// 清理短链接
			for _, link := range links {
				// 清理统计数据
				svcCtx.RepoManager.GetCommonDB().Where("full_short_url = ?", link.FullShortUrl).Delete(&model.LinkAccessLog{})
				svcCtx.RepoManager.GetCommonDB().Where("full_short_url = ?", link.FullShortUrl).Delete(&model.LinkAccessStats{})
				svcCtx.RepoManager.GetCommonDB().Where("full_short_url = ?", link.FullShortUrl).Delete(&model.LinkLocaleStats{})
				svcCtx.RepoManager.GetCommonDB().Where("full_short_url = ?", link.FullShortUrl).Delete(&model.LinkBrowserStats{})
				svcCtx.RepoManager.GetCommonDB().Where("full_short_url = ?", link.FullShortUrl).Delete(&model.LinkOsStats{})
				svcCtx.RepoManager.GetCommonDB().Where("full_short_url = ?", link.FullShortUrl).Delete(&model.LinkDeviceStats{})
				svcCtx.RepoManager.GetCommonDB().Where("full_short_url = ?", link.FullShortUrl).Delete(&model.LinkNetworkStats{})

				// 清理短链接跳转记录
				svcCtx.RepoManager.LinkGoto.Delete(ctx, link.FullShortUrl)

				// 清理短链接
				svcCtx.RepoManager.Link.Delete(ctx, link.ID, link.Gid)
			}
		}

		// 删除测试分组
		svcCtx.RepoManager.Group.DeleteByGidAndUsername(ctx, gid, "test_user")

		// 注意：不删除测试用户，以便复用
	}
}

// TestStatsGetGroup_Success 测试正常获取分组统计数据
func TestStatsGetGroup_Success(t *testing.T) {
	svcCtx, ctx := initGroupTest(t)

	// 准备测试数据
	gid, cleanup := prepareGroupTestData(t, svcCtx, ctx)
	defer cleanup()

	// 测试获取统计数据
	l := logic.NewStatsGetGroupLogic(ctx, svcCtx)
	now := time.Now()
	startDate := now.AddDate(0, 0, -3).Format("2006-01-02") // 3天前
	endDate := now.Format("2006-01-02")                     // 今天

	resp, err := l.StatsGetGroup(&pb.GetGroupStatsRequest{
		Gid:       gid,
		StartDate: startDate,
		EndDate:   endDate,
	})

	if err != nil {
		t.Fatalf("获取分组统计数据失败: %v", err)
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

	t.Logf("成功获取分组统计数据: PV=%d, UV=%d, UIP=%d", resp.Pv, resp.Uv, resp.Uip)
}

// TestStatsGetGroup_InvalidParams 测试无效参数
func TestStatsGetGroup_InvalidParams(t *testing.T) {
	svcCtx, ctx := initGroupTest(t)

	l := logic.NewStatsGetGroupLogic(ctx, svcCtx)

	// 测试空分组ID
	_, err := l.StatsGetGroup(&pb.GetGroupStatsRequest{
		Gid:       "",
		StartDate: "2023-01-01",
		EndDate:   "2023-01-31",
	})

	if err == nil {
		t.Error("期望空分组ID参数校验失败，但实际成功")
	}

	// 测试空日期范围
	_, err = l.StatsGetGroup(&pb.GetGroupStatsRequest{
		Gid:       "test-group",
		StartDate: "",
		EndDate:   "",
	})

	if err == nil {
		t.Error("期望空日期范围参数校验失败，但实际成功")
	}
}

// TestStatsGetGroup_NoData 测试没有数据的情况
func TestStatsGetGroup_NoData(t *testing.T) {
	svcCtx, ctx := initGroupTest(t)

	// 生成一个不存在的分组ID，确保长度不超过32个字符
	nonExistingGid := fmt.Sprintf("non-exist-%d", time.Now().UnixNano()%100000)

	// 先创建这个分组
	testGroup := &model.Group{
		Gid:        nonExistingGid,
		Name:       "不存在数据的测试分组",
		Username:   "test_user", // 分片键
		SortOrder:  1,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		DelFlag:    0,
	}

	err := svcCtx.RepoManager.Group.Create(ctx, testGroup)
	if err != nil {
		t.Fatalf("创建测试分组失败: %v", err)
	}

	// 清理函数
	defer func() {
		svcCtx.RepoManager.GetGroupDB().Where("gid = ? AND username = ?", nonExistingGid, "test_user").Delete(&model.Group{})
	}()

	l := logic.NewStatsGetGroupLogic(ctx, svcCtx)
	resp, err := l.StatsGetGroup(&pb.GetGroupStatsRequest{
		Gid:       nonExistingGid,
		StartDate: "2023-01-01",
		EndDate:   "2023-01-31",
	})

	if err != nil {
		t.Fatalf("获取不存在数据的分组统计数据应该返回空结果而不是错误: %v", err)
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

// TestStatsGetGroup_PermissionDenied 测试权限不足的情况
func TestStatsGetGroup_PermissionDenied(t *testing.T) {
	svcCtx, ctx := initGroupTest(t)

	// 创建一个无权限上下文（没有用户名）
	noAuthCtx := context.Background()

	l := logic.NewStatsGetGroupLogic(noAuthCtx, svcCtx)
	_, err := l.StatsGetGroup(&pb.GetGroupStatsRequest{
		Gid:       "test-group-stats",
		StartDate: "2023-01-01",
		EndDate:   "2023-01-31",
	})

	if err == nil {
		t.Error("期望未登录用户访问时返回错误，但实际成功")
	}

	// 使用一个不存在的分组ID（但用户存在）
	l = logic.NewStatsGetGroupLogic(ctx, svcCtx)
	_, err = l.StatsGetGroup(&pb.GetGroupStatsRequest{
		Gid:       "non-existing-group",
		StartDate: "2023-01-01",
		EndDate:   "2023-01-31",
	})

	if err == nil {
		t.Error("期望不存在的分组ID返回错误，但实际成功")
	}
}
