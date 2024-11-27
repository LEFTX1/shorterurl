package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"shorterurl/api/internal/logic/user"
	"shorterurl/api/internal/svc"
	"shorterurl/api/internal/types"
)

func ApiUserInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserUsernameReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewApiUserInfoLogic(r.Context(), svcCtx)
		resp, err := l.ApiUserInfo(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
