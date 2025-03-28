package logic_test

import (
	"shorterurl/link/rpc/internal/logic"
	"shorterurl/link/rpc/pb"
	"strings"
	"testing"
	"time"
)

// TestShortLinkUpdate_Normal 测试正常更新短链接
func TestShortLinkUpdate_Normal(t *testing.T) {
	// 首先创建一个短链接用于测试
	svcCtx, ctx := setupTest(t)

	// 定义测试用的分组ID
	testGid := "test-update-group"

	// 1. 创建短链接
	createLogic := logic.NewShortLinkCreateLogic(ctx, svcCtx)
	createResp, err := createLogic.ShortLinkCreate(&pb.CreateShortLinkRequest{
		OriginUrl:     "https://github.com/zeromicro/go-zero",
		Domain:        "",
		Gid:           testGid,
		CreatedType:   0,
		ValidDateType: 0,
		ValidDate:     "",
		Describe:      "测试更新前的短链接",
	})

	if err != nil {
		t.Errorf("创建短链接失败: %v", err)
		return
	}

	t.Logf("创建的短链接: %s", createResp.FullShortUrl)

	// 提取短链接URI，去掉http://前缀
	fullShortUrl := createResp.FullShortUrl
	if strings.HasPrefix(fullShortUrl, "http://") {
		fullShortUrl = fullShortUrl[7:]
	}

	t.Logf("处理后的短链接: %s", fullShortUrl)

	// 2. 更新短链接
	updateLogic := logic.NewShortLinkUpdateLogic(ctx, svcCtx)
	_, err = updateLogic.ShortLinkUpdate(&pb.UpdateShortLinkRequest{
		FullShortUrl:  fullShortUrl,
		OriginUrl:     "https://github.com/zeromicro/go-zero/blob/master/readme.md",
		Gid:           testGid, // 使用相同的分组ID
		ValidDateType: 0,
		ValidDate:     "",
		Describe:      "测试更新后的短链接",
	})

	if err != nil {
		t.Errorf("更新短链接失败: %v, 短链接: %s", err, fullShortUrl)
		return
	}

	t.Logf("更新短链接成功: %s", fullShortUrl)
}

// TestShortLinkUpdate_WhitelistEnabled 测试白名单验证
func TestShortLinkUpdate_WhitelistEnabled(t *testing.T) {
	// 首先创建一个短链接用于测试
	svcCtx, ctx := setupTest(t)

	// 定义测试用的分组ID
	testGid := "test-update-whitelist-group"

	// 1. 创建短链接
	createLogic := logic.NewShortLinkCreateLogic(ctx, svcCtx)
	createResp, err := createLogic.ShortLinkCreate(&pb.CreateShortLinkRequest{
		OriginUrl:     "https://github.com/zeromicro/go-zero",
		Domain:        "",
		Gid:           testGid,
		CreatedType:   0,
		ValidDateType: 0,
		ValidDate:     "",
		Describe:      "测试白名单前的短链接",
	})

	if err != nil {
		t.Errorf("创建短链接失败: %v", err)
		return
	}

	t.Logf("创建的短链接: %s", createResp.FullShortUrl)

	// 提取短链接URI，去掉http://前缀
	fullShortUrl := createResp.FullShortUrl
	if strings.HasPrefix(fullShortUrl, "http://") {
		fullShortUrl = fullShortUrl[7:]
	}

	t.Logf("处理后的短链接: %s", fullShortUrl)

	// 2. 尝试更新到不在白名单中的域名
	updateLogic := logic.NewShortLinkUpdateLogic(ctx, svcCtx)
	_, err = updateLogic.ShortLinkUpdate(&pb.UpdateShortLinkRequest{
		FullShortUrl:  fullShortUrl,
		OriginUrl:     "https://example.com/test", // 不在白名单中
		Gid:           testGid,                    // 使用相同的分组ID
		ValidDateType: 0,
		ValidDate:     "",
		Describe:      "测试白名单更新",
	})

	if err == nil {
		t.Error("期望白名单验证失败，但实际成功")
		return
	}

	t.Logf("白名单验证正确拒绝: %v", err)

	// 3. 更新到白名单中的域名
	_, err = updateLogic.ShortLinkUpdate(&pb.UpdateShortLinkRequest{
		FullShortUrl:  fullShortUrl,
		OriginUrl:     "https://baidu.com/search", // 在白名单中
		Gid:           testGid,                    // 使用相同的分组ID
		ValidDateType: 0,
		ValidDate:     "",
		Describe:      "测试白名单更新",
	})

	if err != nil {
		t.Errorf("白名单验证错误拒绝: %v", err)
		return
	}

	t.Logf("白名单验证正确放行: %s", fullShortUrl)
}

// TestShortLinkUpdate_InvalidParams 测试无效参数
func TestShortLinkUpdate_InvalidParams(t *testing.T) {
	svcCtx, ctx := setupTest(t)

	updateLogic := logic.NewShortLinkUpdateLogic(ctx, svcCtx)

	// 测试空短链接
	_, err := updateLogic.ShortLinkUpdate(&pb.UpdateShortLinkRequest{
		FullShortUrl:  "",
		OriginUrl:     "https://github.com/zeromicro/go-zero",
		Gid:           "test-update-invalid",
		ValidDateType: 0,
		ValidDate:     "",
		Describe:      "测试无效参数",
	})

	if err == nil {
		t.Error("期望空短链接验证失败，但实际成功")
		return
	}

	t.Logf("空短链接验证正确拒绝: %v", err)

	// 测试空原始链接
	_, err = updateLogic.ShortLinkUpdate(&pb.UpdateShortLinkRequest{
		FullShortUrl:  "s.xleft.cn/abcdef",
		OriginUrl:     "",
		Gid:           "test-update-invalid",
		ValidDateType: 0,
		ValidDate:     "",
		Describe:      "测试无效参数",
	})

	if err == nil {
		t.Error("期望空原始链接验证失败，但实际成功")
		return
	}

	t.Logf("空原始链接验证正确拒绝: %v", err)

	// 测试空分组标识
	_, err = updateLogic.ShortLinkUpdate(&pb.UpdateShortLinkRequest{
		FullShortUrl:  "s.xleft.cn/abcdef",
		OriginUrl:     "https://github.com/zeromicro/go-zero",
		Gid:           "",
		ValidDateType: 0,
		ValidDate:     "",
		Describe:      "测试无效参数",
	})

	if err == nil {
		t.Error("期望空分组标识验证失败，但实际成功")
		return
	}

	t.Logf("空分组标识验证正确拒绝: %v", err)

	// 测试无效的有效期
	_, err = updateLogic.ShortLinkUpdate(&pb.UpdateShortLinkRequest{
		FullShortUrl:  "s.xleft.cn/abcdef",
		OriginUrl:     "https://github.com/zeromicro/go-zero",
		Gid:           "test-update-invalid",
		ValidDateType: 1,              // 自定义有效期
		ValidDate:     "invalid-date", // 无效日期格式
		Describe:      "测试无效参数",
	})

	if err == nil {
		t.Error("期望无效日期验证失败，但实际成功")
		return
	}

	t.Logf("无效日期验证正确拒绝: %v", err)
}

// TestShortLinkUpdate_WithCustomValidDate 测试自定义有效期
func TestShortLinkUpdate_WithCustomValidDate(t *testing.T) {
	// 首先创建一个短链接用于测试
	svcCtx, ctx := setupTest(t)

	// 定义测试用的分组ID
	testGid := "test-update-date-group"

	// 1. 创建短链接
	createLogic := logic.NewShortLinkCreateLogic(ctx, svcCtx)
	createResp, err := createLogic.ShortLinkCreate(&pb.CreateShortLinkRequest{
		OriginUrl:     "https://github.com/zeromicro/go-zero",
		Domain:        "",
		Gid:           testGid,
		CreatedType:   0,
		ValidDateType: 0,
		ValidDate:     "",
		Describe:      "测试有效期前的短链接",
	})

	if err != nil {
		t.Errorf("创建短链接失败: %v", err)
		return
	}

	t.Logf("创建的短链接: %s", createResp.FullShortUrl)

	// 提取短链接URI，去掉http://前缀
	fullShortUrl := createResp.FullShortUrl
	if strings.HasPrefix(fullShortUrl, "http://") {
		fullShortUrl = fullShortUrl[7:]
	}

	t.Logf("处理后的短链接: %s", fullShortUrl)

	// 2. 更新短链接为自定义有效期
	futureDate := time.Now().AddDate(0, 1, 0) // 一个月后
	validDateStr := futureDate.Format(time.RFC3339)

	updateLogic := logic.NewShortLinkUpdateLogic(ctx, svcCtx)
	_, err = updateLogic.ShortLinkUpdate(&pb.UpdateShortLinkRequest{
		FullShortUrl:  fullShortUrl,
		OriginUrl:     "https://github.com/zeromicro/go-zero",
		Gid:           testGid, // 使用相同的分组ID
		ValidDateType: 1,       // 自定义有效期
		ValidDate:     validDateStr,
		Describe:      "测试有效期后的短链接",
	})

	if err != nil {
		t.Errorf("更新短链接有效期失败: %v", err)
		return
	}

	t.Logf("更新短链接有效期成功: %s, 到期时间: %s", fullShortUrl, validDateStr)
}
