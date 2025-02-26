package file

import (
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest/httpx"
	"main/nicofile/internal/logic/file"
	"main/nicofile/internal/svc"
	"main/nicofile/internal/types"
)

const (
	defaultMultipartMemory = 4 << 20 // 4 MB
)

func UploadChunkHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UploadChunkRequest
		if err := r.ParseMultipartForm(defaultMultipartMemory); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		req.ChunkIndex, _ = strconv.Atoi(r.FormValue("chunkIndex"))
		req.FileName = r.FormValue("filename")
		req.MD5 = r.FormValue("md5")
		req.Ext = r.FormValue("ext")
		//var err error
		File, handler, err := r.FormFile("chunk")
		if err != nil || File == nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := file.NewUploadChunkLogic(r.Context(), svcCtx)
		resp, err := l.UploadChunk(&req, &File, handler)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			//httpx.OkJsonCtx(r.Context(), w, nil)
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
