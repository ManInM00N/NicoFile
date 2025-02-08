package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"main/nicofile/internal/logic/user"
	"main/nicofile/internal/svc"
	"main/nicofile/internal/types"
)

func UserChangePasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.NewPasswordRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewUserChangePasswordLogic(r.Context(), svcCtx)
		resp, err := l.UserChangePassword(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
