package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"shorterurl/api/internal/logic/user"
	"shorterurl/api/internal/svc"
	"shorterurl/api/internal/types"
)

func ApiUpdatePasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserUpdatePasswordReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewApiUpdatePasswordLogic(r.Context(), svcCtx)
		err := l.ApiUpdatePassword(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
