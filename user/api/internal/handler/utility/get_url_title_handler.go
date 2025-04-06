package utility

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"shorterurl/user/api/internal/logic/utility"
	"shorterurl/user/api/internal/svc"
	"shorterurl/user/api/internal/types"
)

// 获取网站标题
func GetUrlTitleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetUrlTitleReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := utility.NewGetUrlTitleLogic(r.Context(), svcCtx)
		resp, err := l.GetUrlTitle(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
