package logic_test

import (
	"net/http"
	"net/http/httptest"
	"shorterurl/link/rpc/internal/logic"
	"shorterurl/link/rpc/pb"
	"testing"
)

// TestUrlTitleGet_Success 测试成功获取URL标题
func TestUrlTitleGet_Success(t *testing.T) {
	// 设置测试环境
	svcCtx, ctx := setupTest(t)

	// 创建测试服务器，返回带有固定标题的HTML页面
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<html><head><title>测试页面标题</title></head><body>测试内容</body></html>`))
	}))
	defer server.Close()

	// 执行测试
	l := logic.NewUrlTitleGetLogic(ctx, svcCtx)
	resp, err := l.UrlTitleGet(&pb.GetUrlTitleRequest{
		Url: server.URL,
	})

	if err != nil {
		t.Fatalf("获取URL标题失败: %v", err)
	}

	if resp == nil {
		t.Fatal("返回结果为空")
	}

	// 验证返回的标题
	expectedTitle := "测试页面标题"
	if resp.Title != expectedTitle {
		t.Errorf("期望获取标题为 %s，实际获取: %s", expectedTitle, resp.Title)
	}

	t.Logf("成功获取URL标题: %s", resp.Title)
}

// TestUrlTitleGet_NoTitle 测试页面无标题
func TestUrlTitleGet_NoTitle(t *testing.T) {
	// 设置测试环境
	svcCtx, ctx := setupTest(t)

	// 创建测试服务器，返回没有标题的HTML页面
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<html><head></head><body>没有标题的页面</body></html>`))
	}))
	defer server.Close()

	// 执行测试
	l := logic.NewUrlTitleGetLogic(ctx, svcCtx)
	resp, err := l.UrlTitleGet(&pb.GetUrlTitleRequest{
		Url: server.URL,
	})

	if err != nil {
		t.Fatalf("获取URL标题失败: %v", err)
	}

	if resp == nil {
		t.Fatal("返回结果为空")
	}

	// 验证返回的标题
	if resp.Title != "未找到页面标题" {
		t.Errorf("期望获取标题为 '未找到页面标题'，实际获取: %s", resp.Title)
	}

	t.Logf("无标题页面处理成功: %s", resp.Title)
}

// TestUrlTitleGet_NotFoundURL 测试无法访问的URL
func TestUrlTitleGet_NotFoundURL(t *testing.T) {
	// 设置测试环境
	svcCtx, ctx := setupTest(t)

	// 创建测试服务器，返回404错误
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	// 执行测试
	l := logic.NewUrlTitleGetLogic(ctx, svcCtx)
	resp, err := l.UrlTitleGet(&pb.GetUrlTitleRequest{
		Url: server.URL,
	})

	if err != nil {
		t.Fatalf("期望返回错误结果而不是错误: %v", err)
	}

	if resp == nil {
		t.Fatal("返回结果为空")
	}

	// 验证返回的标题
	if resp.Title != "无法获取页面标题" {
		t.Errorf("期望获取标题为 '无法获取页面标题'，实际获取: %s", resp.Title)
	}

	t.Logf("404页面处理成功: %s", resp.Title)
}

// TestUrlTitleGet_InvalidURL 测试无效URL
func TestUrlTitleGet_InvalidURL(t *testing.T) {
	// 设置测试环境
	svcCtx, ctx := setupTest(t)

	// 执行测试
	l := logic.NewUrlTitleGetLogic(ctx, svcCtx)
	resp, err := l.UrlTitleGet(&pb.GetUrlTitleRequest{
		Url: "invalid://url",
	})

	if err != nil {
		t.Fatalf("期望返回错误结果而不是错误: %v", err)
	}

	if resp == nil {
		t.Fatal("返回结果为空")
	}

	// 验证返回的标题
	if resp.Title != "无法获取页面标题" {
		t.Errorf("期望获取标题为 '无法获取页面标题'，实际获取: %s", resp.Title)
	}

	t.Logf("无效URL处理成功: %s", resp.Title)
}

// TestUrlTitleGet_EmptyURL 测试空URL
func TestUrlTitleGet_EmptyURL(t *testing.T) {
	// 设置测试环境
	svcCtx, ctx := setupTest(t)

	// 执行测试
	l := logic.NewUrlTitleGetLogic(ctx, svcCtx)
	resp, err := l.UrlTitleGet(&pb.GetUrlTitleRequest{
		Url: "",
	})

	if err == nil {
		t.Error("期望空URL返回错误，但实际没有错误")
		return
	}

	if resp != nil {
		t.Errorf("期望空URL返回空响应，但实际返回: %+v", resp)
	}

	t.Logf("空URL参数校验正确: %v", err)
}
