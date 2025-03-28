package logic_test

import (
	"shorterurl/link/rpc/internal/logic"
	"shorterurl/link/rpc/pb"
	"testing"
)

// TestRestoreUrl_Normal 测试正常情况下的短链接跳转
func TestRestoreUrl_Normal(t *testing.T) {
	// 设置测试环境
	svcCtx, ctx := setupTest(t)

	// 1. 首先创建一个短链接用于测试
	createLogic := logic.NewShortLinkCreateLogic(ctx, svcCtx)
	createResp, err := createLogic.ShortLinkCreate(&pb.CreateShortLinkRequest{
		OriginUrl:     "https://github.com/zeromicro/go-zero",
		Domain:        "",
		Gid:           "test-restore",
		CreatedType:   0,
		ValidDateType: 0,
		ValidDate:     "",
		Describe:      "测试跳转的短链接",
	})

	if err != nil {
		t.Errorf("创建短链接失败: %v", err)
		return
	}

	t.Logf("创建的短链接: %s", createResp.FullShortUrl)

	// 从完整短链接中提取短链接后缀
	shortUri := "Unknown" // 默认值
	if createResp.FullShortUrl != "" {
		// 假设格式为 domain/shortUri
		parts := extractShortUri(createResp.FullShortUrl)
		if parts != "" {
			shortUri = parts
			t.Logf("提取的短链接后缀: %s", shortUri)
		}
	}

	// 2. 测试短链接跳转
	restoreLogic := logic.NewRestoreUrlLogic(ctx, svcCtx)
	restoreResp, err := restoreLogic.RestoreUrl(&pb.RestoreUrlRequest{
		ShortUri: shortUri,
	})

	if err != nil {
		t.Errorf("短链接跳转失败: %v", err)
		return
	}

	// 3. 验证跳转结果
	if restoreResp.OriginUrl != "https://github.com/zeromicro/go-zero" {
		t.Errorf("期望原始链接为 %s, 实际为 %s", "https://github.com/zeromicro/go-zero", restoreResp.OriginUrl)
		return
	}

	t.Logf("短链接跳转成功，原始链接: %s", restoreResp.OriginUrl)
}

// TestRestoreUrl_NotFound 测试不存在的短链接
func TestRestoreUrl_NotFound(t *testing.T) {
	// 设置测试环境
	svcCtx, ctx := setupTest(t)

	// 测试一个不存在的短链接
	restoreLogic := logic.NewRestoreUrlLogic(ctx, svcCtx)
	_, err := restoreLogic.RestoreUrl(&pb.RestoreUrlRequest{
		ShortUri: "not_exist_uri",
	})

	if err == nil {
		t.Error("期望短链接不存在时返回错误，但实际没有错误")
		return
	}

	t.Logf("正确处理不存在的短链接: %v", err)
}

// TestRestoreUrl_EmptyInput 测试空输入
func TestRestoreUrl_EmptyInput(t *testing.T) {
	// 设置测试环境
	svcCtx, ctx := setupTest(t)

	// 测试空短链接
	restoreLogic := logic.NewRestoreUrlLogic(ctx, svcCtx)
	_, err := restoreLogic.RestoreUrl(&pb.RestoreUrlRequest{
		ShortUri: "",
	})

	if err == nil {
		t.Error("期望空短链接时返回错误，但实际没有错误")
		return
	}

	t.Logf("正确处理空短链接: %v", err)
}

// 辅助函数：从完整短链接中提取短链接后缀
func extractShortUri(fullShortUrl string) string {
	if fullShortUrl == "" {
		return ""
	}

	// 如果包含http://或https://，去掉前缀
	if len(fullShortUrl) > 7 && fullShortUrl[:7] == "http://" {
		fullShortUrl = fullShortUrl[7:]
	} else if len(fullShortUrl) > 8 && fullShortUrl[:8] == "https://" {
		fullShortUrl = fullShortUrl[8:]
	}

	// 查找最后一个斜杠的位置
	lastSlashIndex := -1
	for i := len(fullShortUrl) - 1; i >= 0; i-- {
		if fullShortUrl[i] == '/' {
			lastSlashIndex = i
			break
		}
	}

	// 如果找到斜杠，提取后面的部分
	if lastSlashIndex >= 0 && lastSlashIndex < len(fullShortUrl)-1 {
		return fullShortUrl[lastSlashIndex+1:]
	}

	// 没有找到斜杠，返回完整字符串
	return fullShortUrl
}
