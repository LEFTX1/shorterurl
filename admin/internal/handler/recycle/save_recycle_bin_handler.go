package recycle

import (
	"net/http"

	"shorterurl/admin/internal/logic/recycle"
	"shorterurl/admin/internal/svc"
	"shorterurl/admin/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func SaveRecycleBinHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RecycleBinSaveReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := recycle.NewSaveRecycleBinLogic(r.Context(), svcCtx)
		err := l.SaveRecycleBin(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
