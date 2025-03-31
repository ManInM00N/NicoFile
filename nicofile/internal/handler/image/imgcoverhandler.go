package image

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"main/nicofile/internal/logic/image"
	"main/nicofile/internal/svc"
	"main/nicofile/internal/types"
)

func IMGCoverHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DownloadIMGRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := image.NewIMGCoverLogic(r.Context(), svcCtx, w, r)
		err := l.IMGCover(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		}
	}
}
