package logic_test

import (
	"encoding/json"
	"fmt"
	"shorterurl/link/rpc/internal/logic"
	"shorterurl/link/rpc/internal/model"
	"shorterurl/link/rpc/pb"
	"testing"
	"time"
)

// TestShortLinkStats_Redis 测试短链接Redis队列统计
func TestShortLinkStats_Redis(t *testing.T) {
	// 设置测试环境
	svcCtx, ctx := setupTest(t)

	// 模拟短链接数据
	fullShortUrl := "test.example.com/test123"
	gid := "test-group"
	queueKey := "queue:stats:shortlink"

	// 清理队列确保测试环境干净
	_, err := svcCtx.BizRedis.Del(queueKey)
	if err != nil {
		t.Errorf("清理队列失败: %v", err)
		return
	}

	// 测试完成后清理资源
	t.Cleanup(func() {
		_, err := svcCtx.BizRedis.Del(queueKey)
		if err != nil {
			t.Logf("清理测试环境失败: %v", err)
		}
	})

	// 测试短链接统计
	statsLogic := logic.NewShortLinkStatsLogic(ctx, svcCtx)
	_, err = statsLogic.ShortLinkStats(&pb.ShortLinkStatsRequest{
		FullShortUrl: fullShortUrl,
		Gid:          gid,
		User:         "test-user",
		Ip:           "127.0.0.1",
		Browser:      "Chrome",
		Os:           "Windows",
		Device:       "PC",
		Network:      "WIFI",
		Locale:       "CN",
		UvType:       "NEW",
	})

	if err != nil {
		t.Errorf("短链接统计失败: %v", err)
		return
	}

	// 验证统计数据是否被放入Redis队列
	// 检查队列长度
	size, err := svcCtx.BizRedis.Llen(queueKey)
	if err != nil {
		t.Errorf("获取队列长度失败: %v", err)
		return
	}

	// 注意：如果环境中已经有数据，这个判断可能会失败
	if size <= 0 {
		t.Errorf("队列为空，期望至少有一条记录")
		return
	}

	// 获取队列中的最后一条记录（我们刚刚添加的）
	// 注意：这里使用rpop取出最后添加的元素，实际环境中可能会有竞争条件
	record, err := svcCtx.BizRedis.Rpop(queueKey)
	if err != nil {
		t.Errorf("获取队列记录失败: %v", err)
		return
	}

	if record == "" {
		t.Errorf("队列记录为空")
		return
	}

	// 解析JSON记录
	var statsRecord struct {
		FullShortUrl string `json:"fullShortUrl"`
		Gid          string `json:"gid"`
	}

	err = json.Unmarshal([]byte(record), &statsRecord)
	if err != nil {
		t.Errorf("解析记录失败: %v", err)
		return
	}

	// 验证记录数据
	if statsRecord.FullShortUrl != fullShortUrl {
		t.Errorf("短链接不匹配，期望: %s, 实际: %s", fullShortUrl, statsRecord.FullShortUrl)
		return
	}

	if statsRecord.Gid != gid {
		t.Errorf("分组ID不匹配，期望: %s, 实际: %s", gid, statsRecord.Gid)
		return
	}

	// 延时一小段时间确保异步操作有机会尝试执行（但我们不验证结果，因为没有提供RepoManager）
	time.Sleep(time.Millisecond * 100)

	// 清理异步操作产生的Redis缓存
	cacheKey := fmt.Sprintf("link:stats:%s", fullShortUrl)
	_, err = svcCtx.BizRedis.Del(cacheKey)
	if err != nil {
		t.Logf("清理统计缓存失败: %v", err)
	}

	t.Logf("短链接统计测试成功")
}

// TestShortLinkStats_InvalidParams 测试无效参数
func TestShortLinkStats_InvalidParams(t *testing.T) {
	// 设置测试环境
	svcCtx, ctx := setupTest(t)

	statsLogic := logic.NewShortLinkStatsLogic(ctx, svcCtx)

	// 测试空短链接
	_, err := statsLogic.ShortLinkStats(&pb.ShortLinkStatsRequest{
		FullShortUrl: "",
		Gid:          "test-group",
	})

	if err == nil {
		t.Error("期望空短链接时返回错误，但实际没有错误")
		return
	}

	t.Logf("空短链接参数校验正确: %v", err)

	// 测试空分组ID
	_, err = statsLogic.ShortLinkStats(&pb.ShortLinkStatsRequest{
		FullShortUrl: "s.xleft.cn/test",
		Gid:          "",
	})

	if err == nil {
		t.Error("期望空分组ID时返回错误，但实际没有错误")
		return
	}

	t.Logf("空分组ID参数校验正确: %v", err)
}

// TestShortLinkStats_WithDB 使用数据库的完整测试
func TestShortLinkStats_WithDB(t *testing.T) {
	// 如果需要跳过这个测试（例如CI环境中没有数据库），可以使用以下代码
	// t.Skip("跳过需要数据库的测试")

	// 设置测试环境
	svcCtx, ctx := setupTest(t)

	// 清理队列确保测试环境干净
	queueKey := "queue:stats:shortlink"
	_, err := svcCtx.BizRedis.Del(queueKey)
	if err != nil {
		t.Errorf("清理队列失败: %v", err)
		return
	}

	// 准备测试数据
	linkRepo := svcCtx.RepoManager.Link

	// 生成随机的短链接后缀以避免重复，确保长度不超过8字符
	randomID := fmt.Sprintf("t%d", time.Now().UnixNano()%1000000) // 限制在7个字符内 (t + 最多6位数字)

	// 创建测试链接
	testLink := &model.Link{
		Domain:        "test.example.com",
		ShortUri:      randomID,
		FullShortUrl:  fmt.Sprintf("test.example.com/%s", randomID),
		OriginUrl:     "https://github.com/zeromicro/go-zero",
		Gid:           "test-group",
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
		TotalPv:       0,
		TotalUv:       0,
		TotalUip:      0,
		EnableStatus:  0,
		ValidDateType: 0,
		ValidDate:     time.Now().AddDate(10, 0, 0), // 设置有效期为10年后
		DelFlag:       0,
	}

	err = linkRepo.Create(ctx, testLink)
	if err != nil {
		t.Fatalf("创建测试链接失败: %v", err)
	}

	// 创建测试链接跳转记录
	linkGotoRepo := svcCtx.RepoManager.LinkGoto
	testLinkGoto := &model.LinkGoto{
		Gid:          testLink.Gid,
		FullShortUrl: testLink.FullShortUrl,
	}

	err = linkGotoRepo.Create(ctx, testLinkGoto)
	if err != nil {
		t.Fatalf("创建测试链接跳转记录失败: %v", err)
	}

	// 清理函数
	defer func() {
		// 标记链接为已删除
		testLink.DelFlag = 1
		linkRepo.Update(ctx, testLink)

		// 删除链接跳转记录
		linkGotoRepo.Delete(ctx, testLink.FullShortUrl)

		// 清理Redis中的统计数据
		cacheKey := fmt.Sprintf("link:stats:%s", testLink.FullShortUrl)
		svcCtx.BizRedis.Del(cacheKey)

		// 清理Redis队列
		svcCtx.BizRedis.Del(queueKey)
	}()

	// 记录原始PV值
	originalPV := testLink.TotalPv

	// 测试短链接统计
	statsLogic := logic.NewShortLinkStatsLogic(ctx, svcCtx)
	_, err = statsLogic.ShortLinkStats(&pb.ShortLinkStatsRequest{
		FullShortUrl: testLink.FullShortUrl,
		Gid:          testLink.Gid,
		User:         "test-user",
		Ip:           "127.0.0.1",
		Browser:      "Chrome",
		Os:           "Windows",
		Device:       "PC",
		Network:      "WIFI",
		Locale:       "CN",
		UvType:       "NEW",
	})

	if err != nil {
		t.Errorf("短链接统计失败: %v", err)
		return
	}

	// 等待异步更新完成
	time.Sleep(time.Millisecond * 200)

	// 从数据库重新查询链接
	updatedLink, err := linkRepo.FindByFullShortUrlAndGid(ctx, testLink.FullShortUrl, testLink.Gid)
	if err != nil {
		t.Errorf("查询链接失败: %v", err)
		return
	}

	// 验证PV是否更新
	if updatedLink.TotalPv <= originalPV {
		t.Errorf("链接访问量未正确更新，原始值: %d, 更新后: %d", originalPV, updatedLink.TotalPv)
		return
	}

	// 验证Redis中的统计数据
	cacheKey := fmt.Sprintf("link:stats:%s", testLink.FullShortUrl)
	pv, err := svcCtx.BizRedis.Hget(cacheKey, "pv")
	if err != nil {
		t.Logf("从Redis获取PV失败 (可能尚未缓存): %v", err)
	} else if pv == "" {
		t.Logf("Redis中PV未设置 (可能尚未缓存)")
	} else {
		t.Logf("Redis中的PV值: %s", pv)
	}

	t.Logf("数据库测试成功完成")
}
