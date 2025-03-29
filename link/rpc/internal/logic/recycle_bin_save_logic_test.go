package logic_test

import (
	"context"
	"shorterurl/link/rpc/internal/logic"
	"shorterurl/link/rpc/internal/model"
	"shorterurl/link/rpc/internal/svc"
	"shorterurl/link/rpc/pb"
	"testing"
	"time"
)

// 清理特定的测试数据
func cleanSpecificTestData(t *testing.T, svcCtx *svc.ServiceContext, ctx context.Context, fullShortUrl string, gid string) {
	// 尝试通过组合主键查询
	link, err := svcCtx.RepoManager.Link.FindByFullShortUrlAndGid(ctx, fullShortUrl, gid)
	if err == nil && link != nil {
		// 找到了数据，执行删除
		err := svcCtx.RepoManager.Link.Delete(ctx, link.ID, link.Gid)
		if err != nil {
			t.Logf("删除链接数据失败: %v", err)
		} else {
			t.Logf("删除链接ID: %d 成功", link.ID)
		}
	}

	// 直接通过SQL删除任何可能存在的测试数据
	db := svcCtx.DBs.LinkDB
	result := db.Where("full_short_url = ? AND gid = ?", fullShortUrl, gid).Delete(&model.Link{})
	if result.Error != nil {
		t.Logf("SQL删除测试数据失败: %v", result.Error)
	} else if result.RowsAffected > 0 {
		t.Logf("SQL删除 %d 条测试数据", result.RowsAffected)
	}

	t.Logf("已清理特定短链接 %s (分组: %s) 的测试数据", fullShortUrl, gid)
}

// TestRecycleBinSave_Normal 测试正常保存到回收站
func TestRecycleBinSave_Normal(t *testing.T) {
	// 设置测试环境
	svcCtx, ctx := setupTest(t)

	// 创建测试数据
	testGid := "test-recycle-bin-save"
	testFullShortUrl := "test.example.com/sv1"

	// 先清理可能存在的测试数据
	cleanSpecificTestData(t, svcCtx, ctx, testFullShortUrl, testGid)

	// 创建一个测试链接
	testLink := &model.Link{
		Domain:        "test.example.com",
		ShortUri:      "sv1", // 不超过8个字符
		FullShortUrl:  testFullShortUrl,
		OriginUrl:     "https://github.com/zeromicro/go-zero",
		Gid:           testGid,
		EnableStatus:  0, // 启用状态
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
		ValidDateType: 0,
		ValidDate:     time.Now().AddDate(10, 0, 0), // 10年后过期
		DelFlag:       0,
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

	// 测试保存到回收站
	saveLogic := logic.NewRecycleBinSaveLogic(ctx, svcCtx)
	resp, err := saveLogic.RecycleBinSave(&pb.SaveToRecycleBinRequest{
		FullShortUrl: testFullShortUrl,
		Gid:          testGid,
	})

	if err != nil {
		t.Errorf("保存到回收站失败: %v", err)
		return
	}

	if !resp.Success {
		t.Error("保存到回收站响应成功标志为false")
		return
	}

	// 验证链接是否已经进入回收站
	updatedLink, err := svcCtx.RepoManager.Link.FindByFullShortUrlAndGid(ctx, testFullShortUrl, testGid)
	if err != nil {
		t.Errorf("查询更新后的链接失败: %v", err)
		return
	}

	if updatedLink.EnableStatus != 1 {
		t.Errorf("链接未进入回收站，期望EnableStatus=1，实际为%d", updatedLink.EnableStatus)
		return
	}

	t.Logf("正常保存到回收站测试成功")
}

// TestRecycleBinSave_AlreadyInRecycleBin 测试保存已经在回收站中的链接
func TestRecycleBinSave_AlreadyInRecycleBin(t *testing.T) {
	// 设置测试环境
	svcCtx, ctx := setupTest(t)

	// 创建测试数据
	testGid := "test-already-in-recycle-bin"
	testFullShortUrl := "test.example.com/sv2"

	// 先清理可能存在的测试数据
	cleanSpecificTestData(t, svcCtx, ctx, testFullShortUrl, testGid)

	// 创建一个已经在回收站中的测试链接
	testLink := &model.Link{
		Domain:        "test.example.com",
		ShortUri:      "sv2", // 不超过8个字符
		FullShortUrl:  testFullShortUrl,
		OriginUrl:     "https://github.com/zeromicro/go-zero",
		Gid:           testGid,
		EnableStatus:  1, // 已经是未启用状态
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
		ValidDateType: 0,
		ValidDate:     time.Now().AddDate(10, 0, 0), // 10年后过期
		DelFlag:       0,
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

	// 测试保存到回收站
	saveLogic := logic.NewRecycleBinSaveLogic(ctx, svcCtx)
	resp, err := saveLogic.RecycleBinSave(&pb.SaveToRecycleBinRequest{
		FullShortUrl: testFullShortUrl,
		Gid:          testGid,
	})

	if err != nil {
		t.Errorf("保存已在回收站的链接应该成功，但失败: %v", err)
		return
	}

	if !resp.Success {
		t.Error("保存到回收站响应成功标志为false")
		return
	}

	t.Logf("保存已在回收站的链接测试成功")
}

// TestRecycleBinSave_NotExist 测试保存不存在的链接
func TestRecycleBinSave_NotExist(t *testing.T) {
	// 设置测试环境
	svcCtx, ctx := setupTest(t)

	// 测试保存不存在的链接到回收站
	saveLogic := logic.NewRecycleBinSaveLogic(ctx, svcCtx)
	_, err := saveLogic.RecycleBinSave(&pb.SaveToRecycleBinRequest{
		FullShortUrl: "not.exist.url/sv3",
		Gid:          "not-exist-group",
	})

	if err == nil {
		t.Error("保存不存在的链接应该失败，但成功了")
		return
	}

	t.Logf("保存不存在的链接测试成功，正确返回错误: %v", err)
}

// TestRecycleBinSave_InvalidParams 测试无效参数
func TestRecycleBinSave_InvalidParams(t *testing.T) {
	// 设置测试环境
	svcCtx, ctx := setupTest(t)
	saveLogic := logic.NewRecycleBinSaveLogic(ctx, svcCtx)

	// 测试空短链接
	_, err := saveLogic.RecycleBinSave(&pb.SaveToRecycleBinRequest{
		FullShortUrl: "",
		Gid:          "test-group",
	})

	if err == nil {
		t.Error("空短链接参数应该失败，但成功了")
		return
	}

	t.Logf("空短链接参数测试成功，正确返回错误: %v", err)

	// 测试空分组ID
	_, err = saveLogic.RecycleBinSave(&pb.SaveToRecycleBinRequest{
		FullShortUrl: "test.example.com/test",
		Gid:          "",
	})

	if err == nil {
		t.Error("空分组ID参数应该失败，但成功了")
		return
	}

	t.Logf("空分组ID参数测试成功，正确返回错误: %v", err)
}
