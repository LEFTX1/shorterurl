package redirect

import (
	"net/http"

	"shorterurl/user/api/internal/logic/redirect"
	"shorterurl/user/api/internal/svc"
	"shorterurl/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 短链接跳转
func RedirectShortLinkHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ShortLinkRedirectReq
		if err := httpx.Parse(r, &req); err != nil {
			logx.WithContext(r.Context()).Errorf("解析请求参数失败: %v", err)
			http.Error(w, "无效的短链接", http.StatusBadRequest)
			return
		}

		l := redirect.NewRedirectShortLinkLogic(r.Context(), svcCtx)
		err := l.RedirectShortLink(&req, w, r)

		// 错误处理由逻辑层完成，这里不需要额外处理
		if err != nil {
			// 已经在Logic中处理了HTTP响应，这里不再需要额外响应
			logx.WithContext(r.Context()).Errorf("短链接跳转失败: %v", err)
		}
		// 注意：这里不要调用httpx.Ok，因为我们已经在Logic中发送了重定向响应
	}
}
