package logic_test

import (
	"shorterurl/link/rpc/internal/logic"
	"shorterurl/link/rpc/pb"
	"testing"
)

// TestShortLinkBatchCreate_Normal 测试正常批量创建短链接
func TestShortLinkBatchCreate_Normal(t *testing.T) {
	svcCtx, ctx := setupTest(t)

	l := logic.NewShortLinkBatchCreateLogic(ctx, svcCtx)
	resp, err := l.ShortLinkBatchCreate(&pb.BatchCreateShortLinkRequest{
		OriginUrls: []string{
			"https://github.com/zeromicro/go-zero",
			"https://gitee.com/zeromicro/go-zero",
			"https://google.com/search",
		},
		Domain:        "",
		Gid:           "test-batch",
		ValidDateType: 0,
		ValidDate:     "",
		Describe:      "测试批量短链接",
	})

	if err != nil {
		t.Errorf("批量创建短链接失败: %v", err)
		return
	}

	if resp == nil {
		t.Error("返回结果为空")
		return
	}

	// 检查是否正确处理了所有链接
	if len(resp.Results) != 3 {
		t.Errorf("预期生成3个短链接，实际生成%d个", len(resp.Results))
		return
	}

	for i, result := range resp.Results {
		t.Logf("批量创建短链接 %d: %s -> %s", i+1, result.OriginUrl, result.FullShortUrl)
	}
}

// TestShortLinkBatchCreate_WhitelistFiltering 测试白名单过滤
func TestShortLinkBatchCreate_WhitelistFiltering(t *testing.T) {
	svcCtx, ctx := setupTest(t)

	l := logic.NewShortLinkBatchCreateLogic(ctx, svcCtx)
	resp, err := l.ShortLinkBatchCreate(&pb.BatchCreateShortLinkRequest{
		OriginUrls: []string{
			"https://github.com/zeromicro/go-zero", // 白名单允许
			"https://example.com/test",             // 白名单不允许
			"https://baidu.com/search",             // 白名单允许
		},
		Domain:        "",
		Gid:           "test-batch-whitelist",
		ValidDateType: 0,
		ValidDate:     "",
		Describe:      "测试批量白名单过滤",
	})

	if err != nil {
		t.Errorf("批量创建短链接失败: %v", err)
		return
	}

	if resp == nil {
		t.Error("返回结果为空")
		return
	}

	// 检查是否正确过滤了不在白名单中的链接
	if len(resp.Results) != 2 {
		t.Errorf("预期生成2个短链接（白名单过滤），实际生成%d个", len(resp.Results))
		return
	}

	for i, result := range resp.Results {
		t.Logf("批量创建短链接(白名单过滤) %d: %s -> %s", i+1, result.OriginUrl, result.FullShortUrl)
	}
}

// TestShortLinkBatchCreate_EmptyList 测试空列表
func TestShortLinkBatchCreate_EmptyList(t *testing.T) {
	svcCtx, ctx := setupTest(t)

	l := logic.NewShortLinkBatchCreateLogic(ctx, svcCtx)
	_, err := l.ShortLinkBatchCreate(&pb.BatchCreateShortLinkRequest{
		OriginUrls:    []string{},
		Domain:        "",
		Gid:           "test-batch-empty",
		ValidDateType: 0,
		ValidDate:     "",
		Describe:      "测试空列表",
	})

	if err == nil {
		t.Error("期望空列表验证失败，但实际成功")
		return
	}

	t.Logf("空列表验证正确拒绝: %v", err)
}

// TestShortLinkBatchCreate_InvalidDate 测试无效日期
func TestShortLinkBatchCreate_InvalidDate(t *testing.T) {
	svcCtx, ctx := setupTest(t)

	l := logic.NewShortLinkBatchCreateLogic(ctx, svcCtx)
	_, err := l.ShortLinkBatchCreate(&pb.BatchCreateShortLinkRequest{
		OriginUrls: []string{
			"https://github.com/zeromicro/go-zero",
		},
		Domain:        "",
		Gid:           "test-batch-date",
		ValidDateType: 1,              // 自定义有效期
		ValidDate:     "invalid-date", // 无效日期格式
		Describe:      "测试无效日期",
	})

	if err == nil {
		t.Error("期望无效日期验证失败，但实际成功")
		return
	}

	t.Logf("无效日期验证正确拒绝: %v", err)
}
