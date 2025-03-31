package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"main/nicofile/internal/logic/user"
	"main/nicofile/internal/svc"
	"main/nicofile/internal/types"
)

const (
	defaultMultipartMemory = 4 << 20 // 4 MB
)

func UploadCoverHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UploadIMGRequest
		err := r.ParseMultipartForm(defaultMultipartMemory)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		file, header, err := r.FormFile("pic")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			httpx.ErrorCtx(r.Context(), w, err)

			return
		}
		l := user.NewUploadCoverLogic(r.Context(), svcCtx)
		resp, err := l.UploadCover(&req, &file, header)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
