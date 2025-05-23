package file

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"main/nicofile/internal/logic/file"
	"main/nicofile/internal/svc"
	"main/nicofile/internal/types"
)

func MergeChunkHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MergeChunkRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := file.NewMergeChunkLogic(r.Context(), svcCtx)
		resp, err := l.MergeChunk(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
