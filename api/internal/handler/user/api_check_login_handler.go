package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"shorterurl/api/internal/logic/user"
	"shorterurl/api/internal/svc"
)

func ApiCheckLoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewApiCheckLoginLogic(r.Context(), svcCtx)
		resp, err := l.ApiCheckLogin()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
