package redirect

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"shorterurl/link/rpc/shortlinkservice"
	"shorterurl/user/api/internal/middleware"
	"shorterurl/user/api/internal/svc"
	"shorterurl/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// 错误码和消息
const (
	ErrCodeInvalidShortUri = 400
	ErrCodeNotFound        = 404
	ErrCodeForbidden       = 403
	ErrCodeServerError     = 500
)

// 错误消息
const (
	ErrMsgInvalidShortUri = "请提供有效的短链接"
	ErrMsgNotFound        = "您访问的短链接不存在或已被删除"
	ErrMsgForbidden       = "此短链接已过期或被禁用"
	ErrMsgServerError     = "处理您的请求时发生错误，请稍后再试"
)

type RedirectShortLinkLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	response http.ResponseWriter
	request  *http.Request
}

// 短链接跳转
func NewRedirectShortLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RedirectShortLinkLogic {
	return &RedirectShortLinkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 返回美观的错误页面
func (l *RedirectShortLinkLogic) renderErrorPage(w http.ResponseWriter, statusCode int, title, message string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(statusCode)

	// 简单的HTML错误页面模板
	htmlTemplate := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>%s</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background-color: #f5f5f5;
            color: #333;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
            margin: 0;
        }
        .error-container {
            background-color: white;
            border-radius: 8px;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
            padding: 30px;
            text-align: center;
            max-width: 500px;
            width: 100%%;
        }
        h1 {
            color: #e74c3c;
            margin-bottom: 20px;
        }
        p {
            font-size: 18px;
            line-height: 1.6;
        }
        .home-link {
            display: inline-block;
            margin-top: 20px;
            color: #3498db;
            text-decoration: none;
        }
        .home-link:hover {
            text-decoration: underline;
        }
    </style>
</head>
<body>
    <div class="error-container">
        <h1>%s</h1>
        <p>%s</p>
        <a href="/" class="home-link">返回首页</a>
    </div>
</body>
</html>
	`

	html := fmt.Sprintf(htmlTemplate, title, title, message)
	w.Write([]byte(html))
}

// 处理gRPC错误
func (l *RedirectShortLinkLogic) handleGrpcError(err error, w http.ResponseWriter) error {
	grpcStatus, ok := status.FromError(err)
	if ok {
		switch grpcStatus.Code() {
		case codes.NotFound:
			l.renderErrorPage(w, ErrCodeNotFound, "链接不存在", ErrMsgNotFound)
		case codes.PermissionDenied:
			l.renderErrorPage(w, ErrCodeForbidden, "链接已失效", ErrMsgForbidden)
		case codes.InvalidArgument:
			l.renderErrorPage(w, ErrCodeInvalidShortUri, "无效的短链接", ErrMsgInvalidShortUri)
		default:
			l.renderErrorPage(w, ErrCodeServerError, "服务器错误", ErrMsgServerError)
		}
	} else {
		l.renderErrorPage(w, ErrCodeServerError, "服务器错误", ErrMsgServerError)
	}
	return err
}

func (l *RedirectShortLinkLogic) RedirectShortLink(req *types.ShortLinkRedirectReq, w http.ResponseWriter, r *http.Request) error {
	// 保存HTTP响应和请求对象，用于重定向
	l.response = w
	l.request = r

	// 1. 从上下文中获取中间件收集的统计信息
	stats, ok := l.ctx.Value(middleware.RedirectStatsKey).(*types.RedirectStats)
	if !ok {
		l.Logger.Error("未能从上下文中获取统计信息，中间件可能未正确配置")
		l.renderErrorPage(w, ErrCodeServerError, "服务器错误", ErrMsgServerError)
		return nil
	}

	// 2. 参数校验
	if req.ShortUri == "" {
		l.Logger.Error("短链接URI为空")
		l.renderErrorPage(w, ErrCodeInvalidShortUri, "无效的短链接", ErrMsgInvalidShortUri)
		return nil
	}

	// 如果请求的是网站图标，直接返回
	if strings.ToLower(req.ShortUri) == "favicon.ico" {
		w.WriteHeader(http.StatusNotFound)
		return nil
	}

	// 设置短链接URI
	if stats.ShortUri == "" {
		stats.ShortUri = req.ShortUri
	}

	// 3. 调用RPC服务获取原始URL
	// 将统计信息添加到RPC上下文
	ctx := l.ctx

	// 将整个统计信息对象存储到上下文
	ctx = context.WithValue(ctx, middleware.RedirectStatsKey, stats)

	// 将中文地理位置转换为英文
	locale := stats.Locale
	if locale == "本地" {
		locale = "local"
	} else if locale == "未知" {
		locale = "unknown"
	}

	// 使用metadata添加头信息，确保所有值都是ASCII字符
	ctx = metadata.AppendToOutgoingContext(ctx,
		"ip", stats.Ip,
		"user-agent", stats.UserAgent,
		"browser", stats.Browser,
		"os", stats.Os,
		"device", stats.Device,
		"network", "unknown", // 使用英文替代中文
		"locale", locale, // 使用转换后的英文值
	)

	// 记录请求信息
	l.Logger.Infof("短链接跳转请求: URI=%s, IP=%s, UA=%s, Browser=%s, OS=%s, Device=%s, Network=%s, Locale=%s",
		stats.ShortUri, stats.Ip, stats.UserAgent, stats.Browser, stats.Os, stats.Device, stats.Network, stats.Locale)

	// 获取当前用户信息
	userInfo := l.ctx.Value("userInfo").(*types.UserInfo)

	// 创建新的上下文并添加用户信息
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(
		"username", userInfo.Username,
	))

	resp, err := l.svcCtx.LinkRpc.RestoreUrl(ctx, &shortlinkservice.RestoreUrlRequest{
		ShortUri: req.ShortUri,
	})

	// 4. 处理错误情况
	if err != nil {
		l.Logger.Errorf("获取原始URL失败: %v", err)
		return l.handleGrpcError(err, w)
	}

	// 5. 确保原始URL存在
	if resp.OriginUrl == "" {
		l.Logger.Error("获取到的原始URL为空")
		l.renderErrorPage(w, ErrCodeInvalidShortUri, "无效的链接", "此短链接指向的原始URL不存在")
		return nil
	}

	// 6. 记录成功重定向信息
	l.Logger.Infof("短链接 %s 成功重定向到 %s", req.ShortUri, resp.OriginUrl)

	// 7. 执行HTTP重定向
	http.Redirect(w, r, resp.OriginUrl, http.StatusFound)
	return nil
}
