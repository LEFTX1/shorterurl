package group

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"shorterurl/user/api/internal/logic/group"
	"shorterurl/user/api/internal/svc"
	"shorterurl/user/api/internal/types"
)

func CreateGroupHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ShortLinkGroupSaveReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := group.NewCreateGroupLogic(r.Context(), svcCtx)
		resp, err := l.CreateGroup(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}