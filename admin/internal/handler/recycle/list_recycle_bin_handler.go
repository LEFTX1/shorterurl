package recycle

import (
	"net/http"

	"shorterurl/admin/internal/logic/recycle"
	"shorterurl/admin/internal/svc"
	"shorterurl/admin/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListRecycleBinHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RecycleBinPageReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := recycle.NewListRecycleBinLogic(r.Context(), svcCtx)
		resp, err := l.ListRecycleBin(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
