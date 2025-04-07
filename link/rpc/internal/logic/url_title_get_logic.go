package logic

import (
	"context"
	"errors"
	"io"
	"net/http"
	"shorterurl/link/rpc/internal/svc"
	"shorterurl/link/rpc/pb"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/net/html"
)

type UrlTitleGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUrlTitleGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UrlTitleGetLogic {
	return &UrlTitleGetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// --------------------- URL标题功能接口 ---------------------
func (l *UrlTitleGetLogic) UrlTitleGet(in *pb.GetUrlTitleRequest) (*pb.GetUrlTitleResponse, error) {
	// 参数校验
	if in.Url == "" {
		return nil, errors.New("URL不能为空")
	}

	// 确保URL包含协议
	url := in.Url
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	// 创建带超时的HTTP客户端
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		l.Logger.Errorf("创建HTTP请求失败: %v", err)
		return &pb.GetUrlTitleResponse{
			Title: "无法获取页面标题",
		}, nil
	}

	// 设置请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")

	// 最多重试3次
	var resp *http.Response
	for i := 0; i < 3; i++ {
		resp, err = client.Do(req)
		if err == nil {
			break
		}
		l.Logger.Errorf("第%d次获取URL内容失败: %v", i+1, err)
		time.Sleep(time.Second) // 等待1秒后重试
	}

	if err != nil {
		l.Logger.Errorf("获取URL内容失败: %v", err)
		return &pb.GetUrlTitleResponse{
			Title: "无法获取页面标题",
		}, nil
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		l.Logger.Errorf("HTTP请求失败，状态码: %d", resp.StatusCode)
		return &pb.GetUrlTitleResponse{
			Title: "无法获取页面标题",
		}, nil
	}

	// 解析HTML并提取标题
	title, err := l.extractTitle(resp.Body)
	if err != nil {
		l.Logger.Errorf("解析HTML标题失败: %v", err)
		return &pb.GetUrlTitleResponse{
			Title: "无法解析页面标题",
		}, nil
	}

	// 清理标题中的多余空白字符
	title = strings.TrimSpace(title)
	title = strings.ReplaceAll(title, "\n", "")
	title = strings.ReplaceAll(title, "\r", "")
	title = strings.ReplaceAll(title, "\t", "")

	// 如果标题为空，返回默认值
	if title == "" {
		title = "未找到页面标题"
	}

	return &pb.GetUrlTitleResponse{
		Title: title,
	}, nil
}

// extractTitle 从HTML内容中提取标题
func (l *UrlTitleGetLogic) extractTitle(body io.Reader) (string, error) {
	doc, err := html.Parse(body)
	if err != nil {
		return "", err
	}

	var title string
	var crawler func(*html.Node)
	crawler = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
			title = n.FirstChild.Data
			return
		}
		// 如果找不到title标签，尝试查找meta标签中的title
		if n.Type == html.ElementNode && n.Data == "meta" {
			for _, attr := range n.Attr {
				if attr.Key == "property" && attr.Val == "og:title" {
					for _, attr2 := range n.Attr {
						if attr2.Key == "content" {
							title = attr2.Val
							return
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			crawler(c)
		}
	}
	crawler(doc)

	return title, nil
}
