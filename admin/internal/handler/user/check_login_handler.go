package user

import (
	"net/http"

	"shorterurl/admin/internal/logic/user"
	"shorterurl/admin/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func CheckLoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewCheckLoginLogic(r.Context(), svcCtx)
		resp, err := l.CheckLogin()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
