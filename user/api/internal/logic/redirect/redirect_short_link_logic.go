package redirect

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"shorterurl/link/rpc/shortlinkservice"
	"shorterurl/user/api/internal/svc"
	"shorterurl/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (l *RedirectShortLinkLogic) RedirectShortLink(req *types.ShortLinkRedirectReq, w http.ResponseWriter, r *http.Request) error {
	// 保存HTTP响应和请求对象，用于重定向
	l.response = w
	l.request = r

	// 1. 记录请求信息
	l.Logger.Infof("收到短链接跳转请求, ShortUri: %s, IP: %s, UA: %s",
		req.ShortUri, r.RemoteAddr, r.UserAgent())

	// 2. 参数校验
	if req.ShortUri == "" {
		l.Logger.Error("短链接URI为空")
		l.renderErrorPage(w, http.StatusBadRequest, "无效的短链接", "请提供有效的短链接")
		return nil
	}

	// 如果请求的是网站图标，直接返回
	if strings.ToLower(req.ShortUri) == "favicon.ico" {
		w.WriteHeader(http.StatusNotFound)
		return nil
	}

	// 3. 调用RPC服务获取原始URL
	resp, err := l.svcCtx.LinkRpc.RestoreUrl(l.ctx, &shortlinkservice.RestoreUrlRequest{
		ShortUri: req.ShortUri,
	})

	// 4. 处理错误情况
	if err != nil {
		l.Logger.Errorf("获取原始URL失败: %v", err)

		// 处理gRPC错误状态码
		grpcStatus, ok := status.FromError(err)
		if ok {
			switch grpcStatus.Code() {
			case codes.NotFound:
				l.renderErrorPage(w, http.StatusNotFound, "链接不存在", "您访问的短链接不存在或已被删除")
			case codes.PermissionDenied:
				l.renderErrorPage(w, http.StatusForbidden, "链接已失效", "此短链接已过期或被禁用")
			case codes.InvalidArgument:
				l.renderErrorPage(w, http.StatusBadRequest, "无效的短链接", "请检查您的短链接是否正确")
			default:
				l.renderErrorPage(w, http.StatusInternalServerError, "服务器错误", "处理您的请求时发生错误，请稍后再试")
			}
		} else {
			// 非gRPC错误
			l.renderErrorPage(w, http.StatusInternalServerError, "服务器错误", "处理您的请求时发生错误，请稍后再试")
		}
		return err
	}

	// 5. 确保原始URL存在
	if resp.OriginUrl == "" {
		l.Logger.Error("获取到的原始URL为空")
		l.renderErrorPage(w, http.StatusBadRequest, "无效的链接", "此短链接指向的原始URL不存在")
		return nil
	}

	// 6. 记录成功重定向信息
	l.Logger.Infof("短链接 %s 成功重定向到 %s", req.ShortUri, resp.OriginUrl)

	// 7. 执行HTTP重定向
	http.Redirect(w, r, resp.OriginUrl, http.StatusFound)
	return nil
}
