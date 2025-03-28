package logic_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"shorterurl/link/rpc/internal/config"
	"shorterurl/link/rpc/internal/logic"
	"shorterurl/link/rpc/internal/svc"
	"shorterurl/link/rpc/pb"
	"strings"
	"sync"
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	once   sync.Once
	ctx    context.Context
	svcCtx *svc.ServiceContext
)

// setupTest 设置测试环境
func setupTest(t *testing.T) (*svc.ServiceContext, context.Context) {
	once.Do(func() {
		configFile := flag.String("f", "../../etc/link.yaml", "配置文件路径")
		flag.Parse()

		var c config.Config
		conf.MustLoad(*configFile, &c)

		svcCtx = svc.NewServiceContext(c)
		logx.Disable()
		ctx = context.Background()
	})

	return svcCtx, ctx
}

// TestShortLinkCreate_Normal 测试正常创建短链接
func TestShortLinkCreate_Normal(t *testing.T) {
	svcCtx, ctx := setupTest(t)

	l := logic.NewShortLinkCreateLogic(ctx, svcCtx)
	resp, err := l.ShortLinkCreate(&pb.CreateShortLinkRequest{
		OriginUrl:     "https://github.com/zeromicro/go-zero",
		Domain:        "",
		Gid:           "test",
		CreatedType:   0,
		ValidDateType: 0,
		ValidDate:     "",
		Describe:      "测试短链接",
	})

	if err != nil {
		t.Errorf("创建短链接失败: %v", err)
		return
	}

	if resp == nil {
		t.Error("返回结果为空")
		return
	}

	// 记录生成的完整短链接以便后续清理
	fullShortUrl := resp.FullShortUrl
	if fullShortUrl != "" {
		// 从URL中提取域名+路径（去掉http://前缀）
		if strings.HasPrefix(fullShortUrl, "http://") {
			fullShortUrl = fullShortUrl[7:]
		}

		// 清理在测试期间创建的缓存
		t.Cleanup(func() {
			// 清理Redis短链接缓存
			cacheKey := fmt.Sprintf("link:goto:%s", fullShortUrl)
			_, err := svcCtx.BizRedis.Del(cacheKey)
			if err != nil {
				t.Logf("清理短链接缓存失败: %v", err)
			}

			// 清理统计队列缓存
			queueKey := "queue:stats:shortlink"
			_, err = svcCtx.BizRedis.Del(queueKey)
			if err != nil {
				t.Logf("清理统计队列失败: %v", err)
			}

			// 清理可能创建的统计缓存
			statsKey := fmt.Sprintf("link:stats:%s", fullShortUrl)
			_, err = svcCtx.BizRedis.Del(statsKey)
			if err != nil {
				t.Logf("清理统计缓存失败: %v", err)
			}
		})
	}

	t.Logf("创建短链接成功: %s", resp.FullShortUrl)
}

// TestShortLinkCreate_WhitelistEnabled 测试白名单验证
func TestShortLinkCreate_WhitelistEnabled(t *testing.T) {
	svcCtx, ctx := setupTest(t)

	l := logic.NewShortLinkCreateLogic(ctx, svcCtx)
	resp, err := l.ShortLinkCreate(&pb.CreateShortLinkRequest{
		OriginUrl:     "https://example.com/test", // 不在白名单中
		Domain:        "",
		Gid:           "test",
		CreatedType:   0,
		ValidDateType: 0,
		ValidDate:     "",
		Describe:      "测试白名单",
	})

	if err == nil {
		t.Error("期望白名单验证失败，但实际成功")
		return
	}

	t.Logf("白名单验证正确拒绝: %v", err)

	// 测试白名单允许的域名
	resp, err = l.ShortLinkCreate(&pb.CreateShortLinkRequest{
		OriginUrl:     "https://baidu.com/test", // 在白名单中
		Domain:        "",
		Gid:           "test",
		CreatedType:   0,
		ValidDateType: 0,
		ValidDate:     "",
		Describe:      "测试白名单",
	})

	if err != nil {
		t.Errorf("白名单验证错误拒绝: %v", err)
		return
	}

	if resp == nil {
		t.Error("返回结果为空")
		return
	}

	// 记录创建的完整短链接以便后续清理
	fullShortUrl := resp.FullShortUrl
	if fullShortUrl != "" {
		// 从URL中提取域名+路径（去掉http://前缀）
		if strings.HasPrefix(fullShortUrl, "http://") {
			fullShortUrl = fullShortUrl[7:]
		}

		// 清理在测试期间创建的缓存
		t.Cleanup(func() {
			// 清理Redis短链接缓存
			cacheKey := fmt.Sprintf("link:goto:%s", fullShortUrl)
			_, err := svcCtx.BizRedis.Del(cacheKey)
			if err != nil {
				t.Logf("清理短链接缓存失败: %v", err)
			}

			// 清理统计队列缓存
			queueKey := "queue:stats:shortlink"
			_, err = svcCtx.BizRedis.Del(queueKey)
			if err != nil {
				t.Logf("清理统计队列失败: %v", err)
			}

			// 清理可能创建的统计缓存
			statsKey := fmt.Sprintf("link:stats:%s", fullShortUrl)
			_, err = svcCtx.BizRedis.Del(statsKey)
			if err != nil {
				t.Logf("清理统计缓存失败: %v", err)
			}
		})
	}

	t.Logf("白名单验证正确放行: %s", resp.FullShortUrl)
}

// TestShortLinkCreate_InvalidParam 测试无效参数
func TestShortLinkCreate_InvalidParam(t *testing.T) {
	svcCtx, ctx := setupTest(t)

	l := logic.NewShortLinkCreateLogic(ctx, svcCtx)

	// 测试空URL
	_, err := l.ShortLinkCreate(&pb.CreateShortLinkRequest{
		OriginUrl:     "",
		Domain:        "",
		Gid:           "test",
		CreatedType:   0,
		ValidDateType: 0,
		ValidDate:     "",
		Describe:      "测试无效参数",
	})

	if err == nil {
		t.Error("期望原始链接为空验证失败，但实际成功")
		return
	}

	t.Logf("原始链接为空验证正确拒绝: %v", err)

	// 测试无效有效期
	_, err = l.ShortLinkCreate(&pb.CreateShortLinkRequest{
		OriginUrl:     "https://github.com/zeromicro/go-zero",
		Domain:        "",
		Gid:           "test",
		CreatedType:   0,
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

// TestMain 主测试函数
func TestMain(m *testing.M) {
	// 运行测试前的准备
	flag.Parse()
	os.Exit(m.Run())
}
