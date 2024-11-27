package group

import (
	"net/http"

	"shorterurl/admin/internal/logic/group"
	"shorterurl/admin/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
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
