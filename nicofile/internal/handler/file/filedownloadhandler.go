package file

import (
	"errors"
	"github.com/zeromicro/go-zero/rest/httpx"
	"gorm.io/gorm"
	"main/model"
	"main/nicofile/internal/logic/file"
	"main/nicofile/internal/svc"
	"main/nicofile/internal/types"
	"net/http"
)

func FileDownloadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FileDownloadRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		var File model.File
		if err := svcCtx.DB.Model(&model.File{}).Where("file_path = ?", req.Url).First(&File).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			w.Write([]byte("File Not Found"))
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		//w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		w.Header().Set("Surrogate-Control", "no-store")
		w.Header().Set("Content-Disposition", "attachment; filename="+File.FileName)
		w.Header().Set("Content-Transfer-Encoding", "binary")

		//w.Header().Set("Content-Length", strconv.FormatInt(File.Size, 10))
		l := file.NewFileDownloadLogic(r.Context(), svcCtx)

		_, err := l.FileDownload(&req, w, File)
		//f, err := os.OpenFile(svcCtx.Config.StoragePath+"/"+File.FilePath, os.O_RDONLY, 0666)
		//defer f.Close()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			http.ServeFile(w, r, svcCtx.Config.StoragePath+"/"+File.FilePath)
			//httpx.OkJsonCtx(r.Context(), w, resp)
		}

	}
}
