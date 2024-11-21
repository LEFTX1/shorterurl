package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-shorterurl/admin/internal/logic/user"
	"go-zero-shorterurl/admin/internal/svc"
)

func CheckUsernameHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewCheckUsernameLogic(r.Context(), svcCtx)
		resp, err := l.CheckUsername()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
