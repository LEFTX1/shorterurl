package recycle

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"shorterurl/user/api/internal/logic/recycle"
	"shorterurl/user/api/internal/svc"
	"shorterurl/user/api/internal/types"
)

func SaveToRecycleBinHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RecycleBinSaveReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := recycle.NewSaveToRecycleBinLogic(r.Context(), svcCtx)
		resp, err := l.SaveToRecycleBin(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
