package recycle

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"shorterurl/user/api/internal/logic/recycle"
	"shorterurl/user/api/internal/svc"
	"shorterurl/user/api/internal/types"
)

// 移除短链接
func RecycleBinRemoveHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RecycleBinOperateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := recycle.NewRecycleBinRemoveLogic(r.Context(), svcCtx)
		resp, err := l.RecycleBinRemove(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
