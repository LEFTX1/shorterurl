package logic_test

import (
	"context"
	"fmt"
	"shorterurl/link/rpc/internal/logic"
	"shorterurl/link/rpc/internal/model"
	"shorterurl/link/rpc/internal/svc"
	"shorterurl/link/rpc/pb"
	"testing"
	"time"
)

// 测试前清理指定分组的数据
func cleanTestData(t *testing.T, svcCtx *svc.ServiceContext, ctx context.Context, gid string) {
	// 直接通过SQL物理删除测试数据（更彻底的方式）
	db := svcCtx.DBs.LinkDB
	var deletedCount int64

	// 为防止误删，仅删除特定分组的测试数据
	if gid != "" && len(gid) > 5 && (gid[:5] == "test-" || gid == "test") {
		result := db.Where("gid = ?", gid).Delete(&model.Link{})
		if result.Error != nil {
			t.Logf("物理删除测试数据失败 (gid: %s): %v", gid, result.Error)
		} else {
			deletedCount = result.RowsAffected
			t.Logf("物理删除了 %d 条测试数据 (gid: %s)", deletedCount, gid)
		}
	} else {
		t.Logf("跳过清理，无效的测试 gid: %s", gid)
	}
}

// 生成唯一的短链接标识
func generateUniqueShortUri(prefix string, index int) string {
	// 使用前缀和索引生成确定性的、长度可控的标识符
	// 避免使用 time.Now() 来提高确定性和避免冲突
	return fmt.Sprintf("%s%04d", prefix, index) // 例如 "pgn0001", "pgn0002"
}

// TestRecycleBinPage_Normal 测试正常分页查询回收站
func TestRecycleBinPage_Normal(t *testing.T) {
	// 设置测试环境
	svcCtx, ctx := setupTest(t)

	// 固定的测试分组ID
	testGid := "test-recycle-bin-page-normal"

	// 先清理可能存在的测试数据
	cleanTestData(t, svcCtx, ctx, testGid)

	// 创建多个测试链接
	for i := 0; i < 10; i++ {
		shortUri := generateUniqueShortUri("pgn", i) // 生成唯一的短链接标识
		testLink := &model.Link{
			Domain:        "test.example.com",
			ShortUri:      shortUri,
			FullShortUrl:  "test.example.com/" + shortUri,
			OriginUrl:     "https://github.com/zeromicro/go-zero",
			Gid:           testGid,
			EnableStatus:  1, // 未启用状态(在回收站中)
			CreateTime:    time.Now(),
			UpdateTime:    time.Now(),
			ValidDateType: 0,
			ValidDate:     time.Now().AddDate(10, 0, 0), // 10年后过期
			DelFlag:       0,                            // 未删除
		}

		err := svcCtx.RepoManager.Link.Create(ctx, testLink)
		if err != nil {
			t.Fatalf("创建测试链接失败: %v", err)
		}

		// 在测试结束后删除
		defer func(link *model.Link) {
			svcCtx.RepoManager.Link.Delete(ctx, link.ID, link.Gid)
		}(testLink)
	}

	// 测试分页查询
	pageLogic := logic.NewRecycleBinPageLogic(ctx, svcCtx)
	resp, err := pageLogic.RecycleBinPage(&pb.PageRecycleBinShortLinkRequest{
		Gids:    []string{testGid},
		Current: 1,
		Size:    5,
	})

	if err != nil {
		t.Errorf("分页查询回收站失败: %v", err)
		return
	}

	// 验证返回数据
	if resp.Total < 5 {
		t.Errorf("期望总数至少为5，实际为%d", resp.Total)
		return
	}

	if len(resp.Records) != 5 {
		t.Errorf("期望当前页数据量为5，实际为%d", len(resp.Records))
		return
	}

	t.Logf("正常分页查询回收站测试成功，返回数据量: %d, 总数: %d", len(resp.Records), resp.Total)
}

// TestRecycleBinPage_MultiGroup 测试多分组查询回收站
func TestRecycleBinPage_MultiGroup(t *testing.T) {
	// 设置测试环境
	svcCtx, ctx := setupTest(t)

	// 使用固定的测试分组ID
	testGid1 := "test-recycle-bin-group1-multi"
	testGid2 := "test-recycle-bin-group2-multi"

	// 先清理可能存在的测试数据
	cleanTestData(t, svcCtx, ctx, testGid1)
	cleanTestData(t, svcCtx, ctx, testGid2)

	// 为第一个分组创建测试链接
	for i := 0; i < 5; i++ {
		shortUri := generateUniqueShortUri("m", i) // 修改前缀，确保唯一
		testLink := &model.Link{
			Domain:        "test.example.com",
			ShortUri:      shortUri,
			FullShortUrl:  "test.example.com/" + shortUri,
			OriginUrl:     "https://github.com/zeromicro/go-zero",
			Gid:           testGid1,
			EnableStatus:  1, // 未启用状态(在回收站中)
			CreateTime:    time.Now(),
			UpdateTime:    time.Now(),
			ValidDateType: 0,
			ValidDate:     time.Now().AddDate(10, 0, 0), // 10年后过期
			DelFlag:       0,                            // 未删除
		}

		err := svcCtx.RepoManager.Link.Create(ctx, testLink)
		if err != nil {
			t.Fatalf("创建测试链接失败: %v", err)
		}

		// 在测试结束后删除
		defer func(link *model.Link) {
			svcCtx.RepoManager.Link.Delete(ctx, link.ID, link.Gid)
		}(testLink)
	}

	// 为第二个分组创建测试链接
	for i := 0; i < 5; i++ {
		shortUri := generateUniqueShortUri("n", i) // 使用不同前缀
		testLink := &model.Link{
			Domain:        "test.example.com",
			ShortUri:      shortUri,
			FullShortUrl:  "test.example.com/" + shortUri,
			OriginUrl:     "https://github.com/zeromicro/go-zero",
			Gid:           testGid2,
			EnableStatus:  1, // 未启用状态(在回收站中)
			CreateTime:    time.Now(),
			UpdateTime:    time.Now(),
			ValidDateType: 0,
			ValidDate:     time.Now().AddDate(10, 0, 0), // 10年后过期
			DelFlag:       0,                            // 未删除
		}

		err := svcCtx.RepoManager.Link.Create(ctx, testLink)
		if err != nil {
			t.Fatalf("创建测试链接失败: %v", err)
		}

		// 在测试结束后删除
		defer func(link *model.Link) {
			svcCtx.RepoManager.Link.Delete(ctx, link.ID, link.Gid)
		}(testLink)
	}

	// 测试多分组查询
	pageLogic := logic.NewRecycleBinPageLogic(ctx, svcCtx)
	resp, err := pageLogic.RecycleBinPage(&pb.PageRecycleBinShortLinkRequest{
		Gids:    []string{testGid1, testGid2},
		Current: 1,
		Size:    20,
	})

	if err != nil {
		t.Errorf("多分组查询回收站失败: %v", err)
		return
	}

	// 验证返回数据
	if resp.Total < 10 {
		t.Errorf("期望总数至少为10，实际为%d", resp.Total)
		return
	}

	t.Logf("多分组查询回收站测试成功，返回数据量: %d, 总数: %d", len(resp.Records), resp.Total)
}

// TestRecycleBinPage_EmptyResult 测试空结果
func TestRecycleBinPage_EmptyResult(t *testing.T) {
	// 设置测试环境
	svcCtx, ctx := setupTest(t)

	// 测试不存在的分组
	pageLogic := logic.NewRecycleBinPageLogic(ctx, svcCtx)
	resp, err := pageLogic.RecycleBinPage(&pb.PageRecycleBinShortLinkRequest{
		Gids:    []string{"non-existent-group"},
		Current: 1,
		Size:    10,
	})

	if err != nil {
		t.Errorf("查询不存在分组应该成功但返回空结果，但却失败: %v", err)
		return
	}

	if resp.Total != 0 {
		t.Errorf("期望总数为0，实际为%d", resp.Total)
		return
	}

	if len(resp.Records) != 0 {
		t.Errorf("期望数据量为0，实际为%d", len(resp.Records))
		return
	}

	t.Logf("空结果测试成功，返回数据量: %d, 总数: %d", len(resp.Records), resp.Total)
}

// TestRecycleBinPage_InvalidParams 测试无效参数
func TestRecycleBinPage_InvalidParams(t *testing.T) {
	// 设置测试环境
	svcCtx, ctx := setupTest(t)
	pageLogic := logic.NewRecycleBinPageLogic(ctx, svcCtx)

	// 测试空分组ID
	_, err := pageLogic.RecycleBinPage(&pb.PageRecycleBinShortLinkRequest{
		Gids:    []string{},
		Current: 1,
		Size:    10,
	})

	if err == nil {
		t.Error("空分组ID参数应该失败，但成功了")
		return
	}

	t.Logf("空分组ID参数测试成功，正确返回错误: %v", err)

	// 测试无效的当前页
	_, err = pageLogic.RecycleBinPage(&pb.PageRecycleBinShortLinkRequest{
		Gids:    []string{"test-group"},
		Current: 0,
		Size:    10,
	})

	if err != nil {
		t.Errorf("小于等于0的当前页应该自动纠正为1，而不应该返回错误: %v", err)
		return
	}

	// 测试无效的页大小
	_, err = pageLogic.RecycleBinPage(&pb.PageRecycleBinShortLinkRequest{
		Gids:    []string{"test-group"},
		Current: 1,
		Size:    0,
	})

	if err != nil {
		t.Errorf("小于等于0的页大小应该自动纠正为10，而不应该返回错误: %v", err)
		return
	}

	t.Logf("参数自动纠正测试成功")
}

// InitTestData 在测试开始前初始化测试数据
func InitTestData() {
	// 设置测试环境，使用nil作为测试对象
	svcCtx, ctx := setupTest(nil)

	// 清理所有测试组的数据
	testGroups := []string{
		"test-recycle-bin-page-normal",
		"test-recycle-bin-group1-multi",
		"test-recycle-bin-group2-multi",
		"test-recycle-bin-save",
		"test-already-in-recycle-bin",
		"test-recycle-bin-remove",
		"test-not-in-recycle-bin",
		"test-already-deleted",
		"test-recycle-bin-recover-normal",
		"test-already-enabled-normal",
	}

	for _, gid := range testGroups {
		// 查找该分组下的所有链接
		links, err := svcCtx.RepoManager.Link.FindByGidWithCondition(ctx, gid, map[string]interface{}{}, 1, 100)

		if err != nil {
			continue
		}

		// 删除所有找到的链接
		for _, link := range links {
			svcCtx.RepoManager.Link.Delete(ctx, link.ID, link.Gid)
		}
	}
}

// CleanupTestData 在测试结束后清理数据
func CleanupTestData() {
	// 与InitTestData实现相同
	InitTestData()
}
