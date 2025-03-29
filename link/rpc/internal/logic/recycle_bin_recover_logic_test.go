package logic_test

import (
	"shorterurl/link/rpc/internal/logic"
	"shorterurl/link/rpc/internal/model"
	"shorterurl/link/rpc/pb"
	"testing"
	"time"
)

// TestRecycleBinRecover_Normal 测试正常从回收站恢复
func TestRecycleBinRecover_Normal(t *testing.T) {
	// 设置测试环境
	svcCtx, ctx := setupTest(t)

	// 创建测试数据
	testGid := "test-recycle-bin-recover-normal"
	testFullShortUrl := "test.example.com/rtrn1"
	testShortUri := "rtrn1" // 使用一个独特的URI

	// 先清理可能存在的测试数据
	cleanSpecificTestData(t, svcCtx, ctx, testFullShortUrl, testGid)

	// 查找可能存在的同名数据进行删除（先尝试根据全局唯一的full_short_url查找）
	existingLinks, err := svcCtx.RepoManager.Link.FindByGidWithCondition(ctx, testGid, map[string]interface{}{
		"full_short_url": testFullShortUrl,
	}, 1, 10)

	if err == nil && len(existingLinks) > 0 {
		for _, link := range existingLinks {
			t.Logf("软删除已存在的测试数据: %s (ID: %d)", link.FullShortUrl, link.ID)
			// 使用软删除方法
			svcCtx.RepoManager.Link.Delete(ctx, link.ID, link.Gid)
		}
	}

	// 创建一个在回收站中的测试链接
	testLink := &model.Link{
		Domain:        "test.example.com",
		ShortUri:      testShortUri,
		FullShortUrl:  testFullShortUrl,
		OriginUrl:     "https://github.com/zeromicro/go-zero",
		Gid:           testGid,
		EnableStatus:  0, // 回收站中的链接 EnableStatus 仍为 0 (启用状态)
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
		ValidDateType: 0,
		ValidDate:     time.Now().AddDate(10, 0, 0), // 10年后过期
		DelFlag:       1,                            // 标记为删除，表示在回收站中
		DelTime:       time.Now().Unix(),            // 设置删除时间
	}

	// 保存测试链接
	err = svcCtx.RepoManager.Link.Create(ctx, testLink)
	if err != nil {
		t.Fatalf("创建测试链接失败: %v", err)
	}

	// 清理函数
	defer func() {
		// 软删除测试链接
		if testLink.ID > 0 {
			svcCtx.RepoManager.Link.Delete(ctx, testLink.ID, testLink.Gid)
		}
	}()

	// 测试从回收站恢复
	recoverLogic := logic.NewRecycleBinRecoverLogic(ctx, svcCtx)
	resp, err := recoverLogic.RecycleBinRecover(&pb.RecoverFromRecycleBinRequest{
		FullShortUrl: testFullShortUrl,
		Gid:          testGid,
	})

	if err != nil {
		t.Errorf("从回收站恢复失败: %v", err)
		return
	}

	if !resp.Success {
		t.Error("从回收站恢复响应成功标志为false")
		return
	}

	// 验证链接是否已经恢复
	updatedLink, err := svcCtx.RepoManager.Link.FindByFullShortUrlAndGid(ctx, testFullShortUrl, testGid)
	if err != nil {
		t.Errorf("查询更新后的链接失败: %v", err)
		return
	}

	if updatedLink.EnableStatus != 0 {
		t.Errorf("链接未恢复，期望EnableStatus=0，实际为%d", updatedLink.EnableStatus)
		return
	}

	t.Logf("正常从回收站恢复测试成功")
}

// TestRecycleBinRecover_AlreadyEnabled 测试恢复已经是启用状态的链接
func TestRecycleBinRecover_AlreadyEnabled(t *testing.T) {
	// 设置测试环境
	svcCtx, ctx := setupTest(t)

	// 创建测试数据
	testGid := "test-already-enabled-normal"
	testFullShortUrl := "test.example.com/rtae1"
	testShortUri := "rtae1" // 使用一个独特的URI

	// 先清理可能存在的测试数据
	cleanSpecificTestData(t, svcCtx, ctx, testFullShortUrl, testGid)

	// 查找可能存在的同名数据进行删除（先尝试根据全局唯一的full_short_url查找）
	existingLinks, err := svcCtx.RepoManager.Link.FindByGidWithCondition(ctx, testGid, map[string]interface{}{
		"full_short_url": testFullShortUrl,
	}, 1, 10)

	if err == nil && len(existingLinks) > 0 {
		for _, link := range existingLinks {
			t.Logf("软删除已存在的测试数据: %s (ID: %d)", link.FullShortUrl, link.ID)
			// 使用软删除方法
			svcCtx.RepoManager.Link.Delete(ctx, link.ID, link.Gid)
		}
	}

	// 创建一个已经是启用状态的测试链接
	testLink := &model.Link{
		Domain:        "test.example.com",
		ShortUri:      testShortUri,
		FullShortUrl:  testFullShortUrl,
		OriginUrl:     "https://github.com/zeromicro/go-zero",
		Gid:           testGid,
		EnableStatus:  0, // 已经是启用状态
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
		ValidDateType: 0,
		ValidDate:     time.Now().AddDate(10, 0, 0), // 10年后过期
		DelFlag:       0,
	}

	// 保存测试链接
	err = svcCtx.RepoManager.Link.Create(ctx, testLink)
	if err != nil {
		t.Fatalf("创建测试链接失败: %v", err)
	}

	// 清理函数
	defer func() {
		// 软删除测试链接
		if testLink.ID > 0 {
			svcCtx.RepoManager.Link.Delete(ctx, testLink.ID, testLink.Gid)
		}
	}()

	// 测试从回收站恢复
	recoverLogic := logic.NewRecycleBinRecoverLogic(ctx, svcCtx)
	resp, err := recoverLogic.RecycleBinRecover(&pb.RecoverFromRecycleBinRequest{
		FullShortUrl: testFullShortUrl,
		Gid:          testGid,
	})

	if err != nil {
		t.Errorf("恢复已启用的链接应该成功，但失败: %v", err)
		return
	}

	if !resp.Success {
		t.Error("从回收站恢复响应成功标志为false")
		return
	}

	t.Logf("恢复已启用的链接测试成功")
}

// TestRecycleBinRecover_NotExist 测试恢复不存在的链接
func TestRecycleBinRecover_NotExist(t *testing.T) {
	// 设置测试环境
	svcCtx, ctx := setupTest(t)

	// 测试恢复不存在的链接
	recoverLogic := logic.NewRecycleBinRecoverLogic(ctx, svcCtx)
	_, err := recoverLogic.RecycleBinRecover(&pb.RecoverFromRecycleBinRequest{
		FullShortUrl: "not.exist.url/ne1",
		Gid:          "not-exist-group",
	})

	if err == nil {
		t.Error("恢复不存在的链接应该失败，但成功了")
		return
	}

	t.Logf("恢复不存在的链接测试成功，正确返回错误: %v", err)
}

// TestRecycleBinRecover_InvalidParams 测试无效参数
func TestRecycleBinRecover_InvalidParams(t *testing.T) {
	// 设置测试环境
	svcCtx, ctx := setupTest(t)
	recoverLogic := logic.NewRecycleBinRecoverLogic(ctx, svcCtx)

	// 测试空短链接
	_, err := recoverLogic.RecycleBinRecover(&pb.RecoverFromRecycleBinRequest{
		FullShortUrl: "",
		Gid:          "test-group",
	})

	if err == nil {
		t.Error("空短链接参数应该失败，但成功了")
		return
	}

	t.Logf("空短链接参数测试成功，正确返回错误: %v", err)

	// 测试空分组ID
	_, err = recoverLogic.RecycleBinRecover(&pb.RecoverFromRecycleBinRequest{
		FullShortUrl: "test.example.com/test",
		Gid:          "",
	})

	if err == nil {
		t.Error("空分组ID参数应该失败，但成功了")
		return
	}

	t.Logf("空分组ID参数测试成功，正确返回错误: %v", err)
}
