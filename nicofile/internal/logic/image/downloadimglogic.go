package image

import (
	"context"
	"errors"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"main/nicofile/internal/svc"
	"main/nicofile/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DownloadIMGLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	w      http.ResponseWriter
	r      *http.Request
}

func NewDownloadIMGLogic(ctx context.Context, svcCtx *svc.ServiceContext, w http.ResponseWriter, r *http.Request) *DownloadIMGLogic {
	return &DownloadIMGLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		w:      w,
		r:      r,
	}
}

func (l *DownloadIMGLogic) DownloadIMG(req *types.DownloadIMGRequest) (err error) {
	if strings.Contains(req.ID, "..") {
		return errors.New("无效的文件ID")
	}

	filePath := filepath.Join(l.svcCtx.Config.StoragePath, l.svcCtx.Config.IMG.PicPath, req.ID)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return errors.New("文件不存在")
	}

	// 设置正确的Content-Type
	ext := filepath.Ext(filePath)
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	l.w.Header().Set("Content-Type", contentType)
	http.ServeFile(l.w, l.r, filePath)
	return
}
