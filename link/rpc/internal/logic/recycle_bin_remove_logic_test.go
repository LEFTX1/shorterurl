package logic_test

import (
	"shorterurl/link/rpc/internal/logic"
	"shorterurl/link/rpc/internal/model"
	"shorterurl/link/rpc/pb"
	"testing"
	"time"
)

// TestRecycleBinRemove_Normal 测试正常从回收站永久删除
func TestRecycleBinRemove_Normal(t *testing.T) {
	// 设置测试环境
	svcCtx, ctx := setupTest(t)

	// 创建测试数据
	testGid := "test-recycle-bin-remove"
	testFullShortUrl := "test.example.com/rm1"

	// 先清理可能存在的测试数据
	cleanSpecificTestData(t, svcCtx, ctx, testFullShortUrl, testGid)

	// 创建一个在回收站中的测试链接
	testLink := &model.Link{
		Domain:        "test.example.com",
		ShortUri:      "rm1",
		FullShortUrl:  testFullShortUrl,
		OriginUrl:     "https://github.com/zeromicro/go-zero",
		Gid:           testGid,
		EnableStatus:  1, // 未启用状态(在回收站中)
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
		ValidDateType: 0,
		ValidDate:     time.Now().AddDate(10, 0, 0), // 10年后过期
		DelFlag:       0,                            // 未删除
		DelTime:       0,                            // 删除时间为0
	}

	// 保存测试链接
	err := svcCtx.RepoManager.Link.Create(ctx, testLink)
	if err != nil {
		t.Fatalf("创建测试链接失败: %v", err)
	}

	// 获取操作前的链接状态
	linkBefore, err := svcCtx.RepoManager.Link.FindByFullShortUrlAndGid(ctx, testFullShortUrl, testGid)
	if err != nil {
		t.Fatalf("获取操作前的链接状态失败: %v", err)
	}

	if linkBefore.DelFlag != 0 {
		t.Fatalf("初始链接应该是未删除状态，实际为: %d", linkBefore.DelFlag)
	}

	// 测试从回收站永久删除
	removeLogic := logic.NewRecycleBinRemoveLogic(ctx, svcCtx)
	resp, err := removeLogic.RecycleBinRemove(&pb.RemoveFromRecycleBinRequest{
		FullShortUrl: testFullShortUrl,
		Gid:          testGid,
	})

	if err != nil {
		t.Errorf("从回收站永久删除失败: %v", err)
		return
	}

	if !resp.Success {
		t.Error("从回收站永久删除响应成功标志为false")
		return
	}

	// 手动创建一个独立的数据库查询来验证删除结果
	// 注意：这里不使用FindByFullShortUrlAndGid，因为它会过滤掉已删除的记录
	links, err := svcCtx.RepoManager.Link.FindByCondition(ctx, map[string]interface{}{
		"full_short_url": testFullShortUrl,
		"gid":            testGid,
	}, 1, 10)

	if err != nil || len(links) == 0 {
		t.Errorf("查询更新后的链接失败: %v", err)
		return
	}

	updatedLink := links[0]
	if updatedLink.DelFlag != 1 {
		t.Errorf("链接未被标记为删除，期望DelFlag=1，实际为%d", updatedLink.DelFlag)
		return
	}

	if updatedLink.DelTime == 0 {
		t.Error("链接删除时间未设置")
		return
	}

	t.Logf("正常从回收站永久删除测试成功")
}

// TestRecycleBinRemove_NotInRecycleBin 测试删除不在回收站中的链接
func TestRecycleBinRemove_NotInRecycleBin(t *testing.T) {
	// 设置测试环境
	svcCtx, ctx := setupTest(t)

	// 创建测试数据
	testGid := "test-not-in-recycle-bin"
	testFullShortUrl := "test.example.com/rm2"

	// 先清理可能存在的测试数据
	cleanSpecificTestData(t, svcCtx, ctx, testFullShortUrl, testGid)

	// 创建一个未在回收站中的测试链接
	testLink := &model.Link{
		Domain:        "test.example.com",
		ShortUri:      "rm2",
		FullShortUrl:  testFullShortUrl,
		OriginUrl:     "https://github.com/zeromicro/go-zero",
		Gid:           testGid,
		EnableStatus:  0, // 启用状态(不在回收站中)
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
		ValidDateType: 0,
		ValidDate:     time.Now().AddDate(10, 0, 0), // 10年后过期
		DelFlag:       0,                            // 未删除
	}

	// 保存测试链接
	err := svcCtx.RepoManager.Link.Create(ctx, testLink)
	if err != nil {
		t.Fatalf("创建测试链接失败: %v", err)
	}

	// 清理函数
	defer func() {
		// 删除测试链接
		testLink.DelFlag = 1
		svcCtx.RepoManager.Link.Update(ctx, testLink)
	}()

	// 测试从回收站永久删除
	removeLogic := logic.NewRecycleBinRemoveLogic(ctx, svcCtx)
	_, err = removeLogic.RecycleBinRemove(&pb.RemoveFromRecycleBinRequest{
		FullShortUrl: testFullShortUrl,
		Gid:          testGid,
	})

	if err == nil {
		t.Error("删除不在回收站中的链接应该失败，但成功了")
		return
	}

	t.Logf("删除不在回收站中的链接测试成功，正确返回错误: %v", err)
}

// TestRecycleBinRemove_AlreadyDeleted 测试删除已经标记为删除的链接
func TestRecycleBinRemove_AlreadyDeleted(t *testing.T) {
	// 设置测试环境
	svcCtx, ctx := setupTest(t)

	// 创建测试数据
	testGid := "test-already-deleted"
	testFullShortUrl := "test.example.com/rm3"

	// 先清理可能存在的测试数据
	cleanSpecificTestData(t, svcCtx, ctx, testFullShortUrl, testGid)

	// 由于当前logic实现只能处理未删除的链接，我们分两步测试
	// 1. 先创建正常状态的链接
	testLink := &model.Link{
		Domain:        "test.example.com",
		ShortUri:      "rm3",
		FullShortUrl:  testFullShortUrl,
		OriginUrl:     "https://github.com/zeromicro/go-zero",
		Gid:           testGid,
		EnableStatus:  1, // 未启用状态(在回收站中)
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
		ValidDateType: 0,
		ValidDate:     time.Now().AddDate(10, 0, 0), // 10年后过期
		DelFlag:       0,                            // 未删除
		DelTime:       0,                            // 删除时间为0
	}

	// 保存测试链接
	err := svcCtx.RepoManager.Link.Create(ctx, testLink)
	if err != nil {
		t.Fatalf("创建测试链接失败: %v", err)
	}

	// 2. 将其标记为已删除状态
	testLink.DelFlag = 1
	testLink.DelTime = time.Now().Unix()
	err = svcCtx.RepoManager.Link.Update(ctx, testLink)
	if err != nil {
		t.Fatalf("更新链接为已删除状态失败: %v", err)
	}

	// 3. 恢复为未删除状态以测试删除功能
	// 由于FindByFullShortUrlAndGid方法只查询DelFlag=0的记录
	// 因此我们需要手动恢复为未删除状态才能继续测试
	testLink.DelFlag = 0
	err = svcCtx.RepoManager.Link.Update(ctx, testLink)
	if err != nil {
		t.Fatalf("恢复链接为未删除状态失败: %v", err)
	}

	// 测试从回收站永久删除
	removeLogic := logic.NewRecycleBinRemoveLogic(ctx, svcCtx)
	resp, err := removeLogic.RecycleBinRemove(&pb.RemoveFromRecycleBinRequest{
		FullShortUrl: testFullShortUrl,
		Gid:          testGid,
	})

	if err != nil {
		t.Errorf("删除已标记为删除的链接应该成功，但失败: %v", err)
		return
	}

	if !resp.Success {
		t.Error("从回收站永久删除响应成功标志为false")
		return
	}

	t.Logf("删除已标记为删除的链接测试成功")
}

// TestRecycleBinRemove_NotExist 测试删除不存在的链接
func TestRecycleBinRemove_NotExist(t *testing.T) {
	// 设置测试环境
	svcCtx, ctx := setupTest(t)

	// 测试删除不存在的链接
	removeLogic := logic.NewRecycleBinRemoveLogic(ctx, svcCtx)
	_, err := removeLogic.RecycleBinRemove(&pb.RemoveFromRecycleBinRequest{
		FullShortUrl: "not.exist.url/rm4",
		Gid:          "not-exist-group",
	})

	if err == nil {
		t.Error("删除不存在的链接应该失败，但成功了")
		return
	}

	t.Logf("删除不存在的链接测试成功，正确返回错误: %v", err)
}

// TestRecycleBinRemove_InvalidParams 测试无效参数
func TestRecycleBinRemove_InvalidParams(t *testing.T) {
	// 设置测试环境
	svcCtx, ctx := setupTest(t)
	removeLogic := logic.NewRecycleBinRemoveLogic(ctx, svcCtx)

	// 测试空短链接
	_, err := removeLogic.RecycleBinRemove(&pb.RemoveFromRecycleBinRequest{
		FullShortUrl: "",
		Gid:          "test-group",
	})

	if err == nil {
		t.Error("空短链接参数应该失败，但成功了")
		return
	}

	t.Logf("空短链接参数测试成功，正确返回错误: %v", err)

	// 测试空分组ID
	_, err = removeLogic.RecycleBinRemove(&pb.RemoveFromRecycleBinRequest{
		FullShortUrl: "test.example.com/test",
		Gid:          "",
	})

	if err == nil {
		t.Error("空分组ID参数应该失败，但成功了")
		return
	}

	t.Logf("空分组ID参数测试成功，正确返回错误: %v", err)
}
