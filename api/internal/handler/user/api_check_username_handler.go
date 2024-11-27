package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"shorterurl/api/internal/logic/user"
	"shorterurl/api/internal/svc"
	"shorterurl/api/internal/types"
)

func ApiCheckUsernameHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserCheckUsernameReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewApiCheckUsernameLogic(r.Context(), svcCtx)
		resp, err := l.ApiCheckUsername(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
