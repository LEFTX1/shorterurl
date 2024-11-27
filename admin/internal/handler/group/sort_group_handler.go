package group

import (
	"net/http"

	"shorterurl/admin/internal/logic/group"
	"shorterurl/admin/internal/svc"
	"shorterurl/admin/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func SortGroupHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ShortLinkGroupSortReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := group.NewSortGroupLogic(r.Context(), svcCtx)
		err := l.SortGroup(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
