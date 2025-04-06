package stats

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"shorterurl/user/api/internal/logic/stats"
	"shorterurl/user/api/internal/svc"
	"shorterurl/user/api/internal/types"
)

// 分组短链接访问记录查询
func ShortLinkGroupAccessRecordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ShortLinkGroupAccessRecordReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := stats.NewShortLinkGroupAccessRecordLogic(r.Context(), svcCtx)
		resp, err := l.ShortLinkGroupAccessRecord(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
