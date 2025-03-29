package logic

import (
	"context"
	"errors"
	"io"
	"net/http"
	"shorterurl/link/rpc/internal/svc"
	"shorterurl/link/rpc/pb"
	"strings"

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

	// 发送HTTP请求获取页面内容
	resp, err := http.Get(url)
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

	// 返回标题
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
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			crawler(c)
		}
	}
	crawler(doc)

	if title == "" {
		return "未找到页面标题", nil
	}
	return title, nil
}
