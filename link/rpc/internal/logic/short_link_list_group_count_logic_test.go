package logic_test

import (
	"context"
	"shorterurl/link/rpc/internal/logic"
	"shorterurl/link/rpc/internal/svc"
	"shorterurl/link/rpc/pb"
	"testing"
)

// TestShortLinkListGroupCount_Normal 测试正常获取分组短链接数量
func TestShortLinkListGroupCount_Normal(t *testing.T) {
	// 设置测试环境
	svcCtx, ctx := setupTest(t)

	// 创建一些测试数据
	testGids := []string{"test-group-1", "test-group-2", "test-group-3"}

	// 为每个分组创建一些短链接
	createTestShortLinks(t, svcCtx, ctx, testGids)

	// 执行测试
	countLogic := logic.NewShortLinkListGroupCountLogic(ctx, svcCtx)
	resp, err := countLogic.ShortLinkListGroupCount(&pb.GroupShortLinkCountRequest{
		Gids: testGids,
	})

	if err != nil {
		t.Errorf("获取分组短链接数量失败: %v", err)
		return
	}

	// 检查结果
	if len(resp.GroupCounts) == 0 {
		t.Error("期望获取到分组短链接数量，但实际结果为空")
		return
	}

	// 验证每个分组的数量
	for _, groupCount := range resp.GroupCounts {
		t.Logf("分组 %s 的短链接数量: %d", groupCount.Gid, groupCount.ShortLinkCount)

		// 检查分组ID是否在请求列表中
		found := false
		for _, gid := range testGids {
			if gid == groupCount.Gid {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("返回了未请求的分组ID: %s", groupCount.Gid)
		}
	}
}

// TestShortLinkListGroupCount_EmptyInput 测试空输入
func TestShortLinkListGroupCount_EmptyInput(t *testing.T) {
	// 设置测试环境
	svcCtx, ctx := setupTest(t)

	// 执行测试
	countLogic := logic.NewShortLinkListGroupCountLogic(ctx, svcCtx)
	_, err := countLogic.ShortLinkListGroupCount(&pb.GroupShortLinkCountRequest{
		Gids: []string{},
	})

	if err == nil {
		t.Error("期望空输入时返回错误，但实际没有错误")
		return
	}

	t.Logf("空输入正确返回错误: %v", err)
}

// createTestShortLinks 创建测试短链接数据
func createTestShortLinks(t *testing.T, svcCtx *svc.ServiceContext, ctx context.Context, gids []string) {
	createLogic := logic.NewShortLinkCreateLogic(ctx, svcCtx)

	// 为每个分组创建不同数量的短链接
	for i, gid := range gids {
		// 为每个分组创建 i+1 个短链接
		for j := 0; j < i+1; j++ {
			_, err := createLogic.ShortLinkCreate(&pb.CreateShortLinkRequest{
				OriginUrl:     "https://github.com/zeromicro/go-zero",
				Domain:        "",
				Gid:           gid,
				CreatedType:   0,
				ValidDateType: 0,
				ValidDate:     "",
				Describe:      "测试短链接",
			})

			if err != nil {
				t.Logf("创建测试短链接失败，分组: %s, 错误: %v", gid, err)
			}
		}
	}
}
