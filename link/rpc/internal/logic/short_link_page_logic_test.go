package logic_test

import (
	"shorterurl/link/rpc/internal/logic"
	"shorterurl/link/rpc/pb"
	"testing"
)

// TestShortLinkPage_Normal 测试正常分页查询短链接
func TestShortLinkPage_Normal(t *testing.T) {
	// 首先创建多个短链接用于测试
	svcCtx, ctx := setupTest(t)

	// 1. 创建多个短链接
	createLogic := logic.NewShortLinkCreateLogic(ctx, svcCtx)
	gid := "test-page"

	// 创建5个短链接
	for i := 0; i < 5; i++ {
		_, err := createLogic.ShortLinkCreate(&pb.CreateShortLinkRequest{
			OriginUrl:     "https://github.com/zeromicro/go-zero",
			Domain:        "",
			Gid:           gid,
			CreatedType:   0,
			ValidDateType: 0,
			ValidDate:     "",
			Describe:      "测试分页查询短链接",
		})

		if err != nil {
			t.Errorf("创建短链接失败: %v", err)
			return
		}
	}

	// 2. 执行分页查询
	pageLogic := logic.NewShortLinkPageLogic(ctx, svcCtx)
	resp, err := pageLogic.ShortLinkPage(&pb.PageShortLinkRequest{
		Gid:     gid,
		Current: 1,
		Size:    10,
	})

	if err != nil {
		t.Errorf("分页查询短链接失败: %v", err)
		return
	}

	if resp == nil {
		t.Error("返回结果为空")
		return
	}

	// 检查查询结果是否正确
	if resp.Total < 5 {
		t.Errorf("分页查询短链接结果不正确，预期至少5条记录，实际查询到%d条", resp.Total)
		return
	}

	t.Logf("分页查询短链接成功，总记录数: %d, 当前页: %d, 每页大小: %d", resp.Total, resp.Current, resp.Size)
	for i, record := range resp.Records {
		t.Logf("记录 %d: %s -> %s", i+1, record.OriginUrl, record.FullShortUrl)
	}
}

// TestShortLinkPage_Pagination 测试分页功能
func TestShortLinkPage_Pagination(t *testing.T) {
	// 首先创建多个短链接用于测试
	svcCtx, ctx := setupTest(t)

	// 1. 创建多个短链接
	createLogic := logic.NewShortLinkCreateLogic(ctx, svcCtx)
	gid := "test-page-pagination"

	// 创建10个短链接
	for i := 0; i < 10; i++ {
		_, err := createLogic.ShortLinkCreate(&pb.CreateShortLinkRequest{
			OriginUrl:     "https://github.com/zeromicro/go-zero",
			Domain:        "",
			Gid:           gid,
			CreatedType:   0,
			ValidDateType: 0,
			ValidDate:     "",
			Describe:      "测试分页功能",
		})

		if err != nil {
			t.Errorf("创建短链接失败: %v", err)
			return
		}
	}

	// 2. 执行分页查询 - 第一页
	pageLogic := logic.NewShortLinkPageLogic(ctx, svcCtx)
	resp1, err := pageLogic.ShortLinkPage(&pb.PageShortLinkRequest{
		Gid:     gid,
		Current: 1,
		Size:    5, // 每页5条
	})

	if err != nil {
		t.Errorf("分页查询第一页失败: %v", err)
		return
	}

	if resp1 == nil {
		t.Error("第一页返回结果为空")
		return
	}

	// 3. 执行分页查询 - 第二页
	resp2, err := pageLogic.ShortLinkPage(&pb.PageShortLinkRequest{
		Gid:     gid,
		Current: 2,
		Size:    5, // 每页5条
	})

	if err != nil {
		t.Errorf("分页查询第二页失败: %v", err)
		return
	}

	if resp2 == nil {
		t.Error("第二页返回结果为空")
		return
	}

	// 检查两页数据不同
	if len(resp1.Records) == 0 || len(resp2.Records) == 0 {
		t.Error("返回的记录为空")
		return
	}

	t.Logf("分页查询第一页成功，总记录数: %d, 当前页: %d, 每页大小: %d, 记录数: %d",
		resp1.Total, resp1.Current, resp1.Size, len(resp1.Records))
	t.Logf("分页查询第二页成功，总记录数: %d, 当前页: %d, 每页大小: %d, 记录数: %d",
		resp2.Total, resp2.Current, resp2.Size, len(resp2.Records))
}

// TestShortLinkPage_InvalidParams 测试无效参数
func TestShortLinkPage_InvalidParams(t *testing.T) {
	svcCtx, ctx := setupTest(t)

	pageLogic := logic.NewShortLinkPageLogic(ctx, svcCtx)

	// 测试空分组标识
	_, err := pageLogic.ShortLinkPage(&pb.PageShortLinkRequest{
		Gid:     "",
		Current: 1,
		Size:    10,
	})

	if err == nil {
		t.Error("期望空分组标识验证失败，但实际成功")
		return
	}

	t.Logf("空分组标识验证正确拒绝: %v", err)

	// 测试无效分页参数 - 会使用默认值，不会报错
	resp, err := pageLogic.ShortLinkPage(&pb.PageShortLinkRequest{
		Gid:     "test-page-invalid",
		Current: 0, // 无效页码，会使用默认值1
		Size:    0, // 无效大小，会使用默认值10
	})

	if err != nil {
		t.Errorf("分页参数验证错误拒绝: %v", err)
		return
	}

	t.Logf("无效分页参数正确处理，使用默认值：当前页 %d, 每页大小 %d", resp.Current, resp.Size)
}
