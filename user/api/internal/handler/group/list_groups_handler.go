package group

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"shorterurl/user/api/internal/logic/group"
	"shorterurl/user/api/internal/svc"
)

func ListGroupsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := group.NewListGroupsLogic(r.Context(), svcCtx)
		resp, err := l.ListGroups()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
