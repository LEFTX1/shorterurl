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

// prepareGroupAccessLogTestData 为测试准备分组访问日志数据
func prepareGroupAccessLogTestData(t *testing.T, svcCtx *svc.ServiceContext, ctx context.Context) (string, func()) {
	// 生成唯一的分组ID
	gid := fmt.Sprintf("test-group-acc-%d", time.Now().UnixNano()%100000)

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
		Name:       "测试分组记录查询",
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
	var linkIDs []int64
	var fullShortUrls []string

	for i := 0; i < 3; i++ {
		// 增加随机性，避免重复
		randomID := fmt.Sprintf("gr%d%d", time.Now().UnixNano()%10000, i)
		fullShortUrl := fmt.Sprintf("test.example.com/%s", randomID)
		fullShortUrls = append(fullShortUrls, fullShortUrl)

		testLink := &model.Link{
			Domain:        "test.example.com",
			ShortUri:      randomID,
			FullShortUrl:  fullShortUrl,
			OriginUrl:     fmt.Sprintf("https://github.com/zeromicro/go-zero/page%d", i),
			Gid:           gid, // 分片键
			Favicon:       "",
			EnableStatus:  0,
			CreatedType:   0,
			ValidDateType: 0,
			ValidDate:     time.Now().AddDate(1, 0, 0),
			Describe:      fmt.Sprintf("测试分组访问记录链接-%d", i),
			CreateTime:    time.Now(),
			UpdateTime:    time.Now(),
			DelFlag:       0,
		}

		err = svcCtx.RepoManager.Link.Create(ctx, testLink)
		if err != nil {
			t.Fatalf("创建测试短链接失败: %v", err)
		}
		linkIDs = append(linkIDs, testLink.ID)

		// 创建短链接跳转记录
		testLinkGoto := &model.LinkGoto{
			Gid:          gid,
			FullShortUrl: fullShortUrl,
		}

		err = svcCtx.RepoManager.LinkGoto.Create(ctx, testLinkGoto)
		if err != nil {
			t.Fatalf("创建测试短链接跳转记录失败: %v", err)
		}

		// 为每个链接创建访问日志
		for j := 0; j < 3; j++ {
			accessLog := &model.LinkAccessLog{
				FullShortUrl: fullShortUrl,
				User:         fmt.Sprintf("user%d_%d", i, j),
				IP:           fmt.Sprintf("192.168.%d.%d", i, j),
				Browser:      "Chrome",
				Os:           "Windows",
				Network:      "WIFI",
				Device:       "PC",
				Locale:       "北京",
				CreateTime:   time.Now().Add(-time.Duration(j) * time.Hour),
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
	cleanup := func() {
		// 删除创建的测试数据
		for _, fullShortUrl := range fullShortUrls {
			svcCtx.RepoManager.GetCommonDB().Where("full_short_url = ?", fullShortUrl).Delete(&model.LinkAccessLog{})
			// 使用 LinkGotoRepo 删除，并传入 gid 和 fullShortUrl
			err := svcCtx.RepoManager.LinkGoto.DeleteByGidAndFullShortUrl(ctx, gid, fullShortUrl)
			if err != nil {
				t.Logf("清理LinkGoto记录失败 (gid=%s, url=%s): %v", gid, fullShortUrl, err)
			}
		}

		for _, id := range linkIDs {
			svcCtx.RepoManager.Link.Delete(ctx, id, gid)
		}

		// 使用 GroupRepo 删除分组，并传入 gid
		err := svcCtx.RepoManager.Group.DeleteByGid(ctx, gid)
		if err != nil {
			t.Logf("清理测试分组失败 (gid=%s): %v", gid, err)
		}

		// 清理测试用户
		err = svcCtx.RepoManager.User.DeleteByUsername(ctx, "test_user")
		if err != nil {
			t.Logf("清理测试用户失败: %v", err)
		}

		t.Logf("清理分组访问记录测试数据完成")
	}

	return gid, cleanup
}

// TestStatsGroupAccessRecordQuery_Success 测试成功查询分组访问记录
func TestStatsGroupAccessRecordQuery_Success(t *testing.T) {
	// 设置测试环境
	svcCtx, ctx := initTest(t)

	// 准备测试数据
	gid, cleanup := prepareGroupAccessLogTestData(t, svcCtx, ctx)
	defer cleanup()

	// 创建当前日期
	now := time.Now()
	startDate := now.AddDate(0, 0, -7).Format("2006-01-02")
	endDate := now.Format("2006-01-02")

	// 测试查询分组访问记录
	groupAccessLogic := logic.NewStatsGroupAccessRecordQueryLogic(ctx, svcCtx)
	resp, err := groupAccessLogic.StatsGroupAccessRecordQuery(&pb.GroupAccessRecordQueryRequest{
		Gid:       gid,
		StartDate: startDate,
		EndDate:   endDate,
		Current:   1,
		Size:      10,
	})

	if err != nil {
		t.Errorf("查询分组访问记录失败: %v", err)
		return
	}

	// 验证返回结果
	if resp == nil {
		t.Error("返回结果为空")
		return
	}

	if resp.Total < 1 {
		t.Errorf("期望至少有1条记录，实际有%d条", resp.Total)
		return
	}

	t.Logf("成功查询分组访问记录: 总记录数 %d", resp.Total)
}

// TestStatsGroupAccessRecordQuery_InvalidParams 测试无效参数
func TestStatsGroupAccessRecordQuery_InvalidParams(t *testing.T) {
	// 设置测试环境
	svcCtx, ctx := initTest(t)

	now := time.Now()
	startDate := now.AddDate(0, 0, -7).Format("2006-01-02")
	endDate := now.Format("2006-01-02")

	// 创建逻辑实例
	groupAccessLogic := logic.NewStatsGroupAccessRecordQueryLogic(ctx, svcCtx)

	// 测试空分组ID
	_, err := groupAccessLogic.StatsGroupAccessRecordQuery(&pb.GroupAccessRecordQueryRequest{
		Gid:       "",
		StartDate: startDate,
		EndDate:   endDate,
		Current:   1,
		Size:      10,
	})

	if err == nil {
		t.Error("空分组ID参数校验应该失败")
		return
	}
	t.Logf("空分组ID参数校验正确: %v", err)

	// 测试空日期
	_, err = groupAccessLogic.StatsGroupAccessRecordQuery(&pb.GroupAccessRecordQueryRequest{
		Gid:       "test-gid",
		StartDate: "",
		EndDate:   "",
		Current:   1,
		Size:      10,
	})

	if err == nil {
		t.Error("空日期参数校验应该失败")
		return
	}
	t.Logf("空日期参数校验正确: %v", err)

	// 测试无效分页参数
	_, err = groupAccessLogic.StatsGroupAccessRecordQuery(&pb.GroupAccessRecordQueryRequest{
		Gid:       "test-gid",
		StartDate: startDate,
		EndDate:   endDate,
		Current:   0,
		Size:      10,
	})

	if err == nil {
		t.Error("无效页码参数校验应该失败")
		return
	}
	t.Logf("无效页码参数校验正确: %v", err)

	_, err = groupAccessLogic.StatsGroupAccessRecordQuery(&pb.GroupAccessRecordQueryRequest{
		Gid:       "test-gid",
		StartDate: startDate,
		EndDate:   endDate,
		Current:   1,
		Size:      0,
	})

	if err == nil {
		t.Error("无效每页数量参数校验应该失败")
		return
	}
	t.Logf("无效每页数量参数校验正确: %v", err)
}

// TestStatsGroupAccessRecordQuery_NoData 测试无数据情况
func TestStatsGroupAccessRecordQuery_NoData(t *testing.T) {
	// 设置测试环境
	svcCtx, ctx := initTest(t)

	// 使用一个不存在短链接的分组ID
	gid := fmt.Sprintf("test-no-data-%d", time.Now().UnixNano())

	now := time.Now()
	startDate := now.AddDate(0, 0, -7).Format("2006-01-02")
	endDate := now.Format("2006-01-02")

	// 创建测试分组以通过权限验证
	testGroup := &model.Group{
		Gid:        gid,
		Name:       "测试分组NoData",
		Username:   "test_user",
		SortOrder:  1,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		DelFlag:    0,
	}

	err := svcCtx.RepoManager.Group.Create(ctx, testGroup)
	if err != nil {
		t.Fatalf("创建测试分组失败: %v", err)
	}

	defer func() {
		// 清理测试分组
		err := svcCtx.RepoManager.Group.DeleteByGidAndUsername(ctx, gid, "test_user")
		if err != nil {
			t.Logf("清理测试分组失败: %v", err)
		} else {
			t.Logf("成功清理测试分组")
		}
	}()

	// 执行查询
	groupAccessLogic := logic.NewStatsGroupAccessRecordQueryLogic(ctx, svcCtx)
	resp, err := groupAccessLogic.StatsGroupAccessRecordQuery(&pb.GroupAccessRecordQueryRequest{
		Gid:       gid,
		StartDate: startDate,
		EndDate:   endDate,
		Current:   1,
		Size:      10,
	})

	if err != nil {
		t.Errorf("应该成功返回空结果，但失败: %v", err)
		return
	}

	// 验证返回空结果
	if resp.Total != 0 || len(resp.Records) != 0 {
		t.Errorf("期望返回0条记录，但返回了%d条", resp.Total)
		return
	}

	t.Logf("正确处理无数据情况: 返回0条记录")
}

// TestStatsGroupAccessRecordQuery_PermissionDenied 测试权限拒绝
func TestStatsGroupAccessRecordQuery_PermissionDenied(t *testing.T) {
	// 设置测试环境
	svcCtx, ctx := initTest(t)

	// 创建一个与当前用户不匹配的分组ID
	gid := "test-other-user-group"

	now := time.Now()
	startDate := now.AddDate(0, 0, -7).Format("2006-01-02")
	endDate := now.Format("2006-01-02")

	// 使用错误的用户上下文
	mockCtx := context.WithValue(ctx, "username", "wrong_user")

	// 执行查询
	groupAccessLogic := logic.NewStatsGroupAccessRecordQueryLogic(mockCtx, svcCtx)
	_, err := groupAccessLogic.StatsGroupAccessRecordQuery(&pb.GroupAccessRecordQueryRequest{
		Gid:       gid,
		StartDate: startDate,
		EndDate:   endDate,
		Current:   1,
		Size:      10,
	})

	// 期望返回权限拒绝错误
	if err == nil {
		t.Error("应该返回权限拒绝错误，但没有返回错误")
		return
	}

	t.Logf("正确处理权限拒绝情况: %v", err)
}
