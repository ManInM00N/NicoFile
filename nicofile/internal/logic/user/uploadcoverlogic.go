package user

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"io/ioutil"
	"main/model"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"main/nicofile/internal/svc"
	"main/nicofile/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadCoverLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadCoverLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadCoverLogic {
	return &UploadCoverLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadCoverLogic) UploadCover(req *types.UploadIMGRequest, file *multipart.File, header *multipart.FileHeader) (resp *types.UploadIMGResponse, err error) {
	// 检查文件类型
	buff := make([]byte, 512)
	_, err2 := (*file).Read(buff)
	if err2 != nil {
		return nil, errors.New("文件非法")
	}

	storagePath := filepath.Join(l.svcCtx.Config.StoragePath, l.svcCtx.Config.IMG.PicPath)
	contentType := http.DetectContentType(buff)
	if !isAllowedType(contentType, l.svcCtx.Config.IMG.AllowedTypes) {
		return nil, errors.New("不允许的文件类型")
	}

	// 重置文件指针
	_, err = (*file).Seek(0, 0)
	if err != nil {
		return nil, err
	}
	// 生成唯一文件名
	ext := filepath.Ext(header.Filename)
	newFilename := generateFilename() + ext
	savePath := filepath.Join(storagePath, newFilename)

	// 确保目录存在
	if err2 := os.MkdirAll(storagePath, 0755); err2 != nil {
		return nil, errors.New("目录生成失败")
	}

	// 保存文件
	fileBytes, err := ioutil.ReadAll(*file)
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(savePath, fileBytes, 0644); err != nil {
		return nil, err
	}
	id, _ := l.ctx.Value("UserId").(json.Number).Int64()
	if err2 := l.svcCtx.DB.Model(&model.Image{}).Create(&model.Image{Name: header.Filename, Path: newFilename, AuthorID: (uint)(id)}).Error; err2 != nil {
		return nil, errors.New("数据库错误")
	}
	if err2 := l.svcCtx.DB.Model(&model.User{}).Where("id = ?", id).Update("cover", newFilename).Error; err2 != nil {
		return nil, errors.New("数据库错误")
	}
	return &types.UploadIMGResponse{
		URL: newFilename,
	}, nil
	return
}
func isAllowedType(contentType string, allowedTypes []string) bool {
	for _, t := range allowedTypes {
		if strings.EqualFold(contentType, t) {
			return true
		}
	}
	return false
}

func generateFilename() string {
	return time.Now().Format("20060102150405") + "-" + uuid.New().String()
}
