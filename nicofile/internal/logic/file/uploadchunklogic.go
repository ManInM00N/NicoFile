package file

import (
	"bufio"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"main/model"
	"main/nicofile/internal/svc"
	"main/nicofile/internal/types"
	"main/pkg/util"
	"mime/multipart"
	"os"
	"path/filepath"
)

type UploadChunkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadChunkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadChunkLogic {
	return &UploadChunkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadChunkLogic) UploadChunk(req *types.UploadChunkRequest, File *multipart.File, Handler *multipart.FileHeader) (resp *types.UploadChunkResponse, err error) {
	num := int64(0)
	resp = &types.UploadChunkResponse{
		Error: false,
	}
	id, _ := l.ctx.Value("UserId").(json.Number).Int64()
	var chunk model.Chunk
	l.svcCtx.DB.Model(&model.Chunk{}).Select("id,md5,file_path").Where("file_name = ? and chunk_index = ? and author_id = ? and ext = ?", req.FileName, req.ChunkIndex, id, req.Ext).First(&chunk).Count(&num)
	storagePath := fmt.Sprintf("%s/%d", l.svcCtx.Config.ChunkStorePath, id)
	filePath := fmt.Sprintf("%s/%s_%d", storagePath, req.FileName, req.ChunkIndex)
	if num >= 1 {
		if chunk.MD5 == req.MD5 {
			resp.Message = "文件已存在"
			return
		} else {
			os.Remove(filepath.Join(l.svcCtx.Config.ChunkStorePath, chunk.FilePath))
			l.svcCtx.DB.Unscoped().Model(&model.File{}).Delete(&chunk)
		}
	}
	os.MkdirAll(storagePath, os.ModePerm)
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		resp.Error = true
		resp.Message = "创建分片文件失败"
		util.Log.Errorf("create chunk file error: %v\n", err)
		return resp, nil
	}
	writer := bufio.NewWriter(f)
	arr := util.Pool.Get().([]byte)
	defer util.Pool.Put(arr)
	hash := md5.New()
	defer (*File).Close()
	for {
		len, err2 := (*File).Read(arr)
		if err2 == io.EOF || len == 0 {
			break
		}
		if err2 != nil && err2 != io.EOF {
			resp.Error = true
			resp.Message = "上传失败"
			f.Close()
			os.Remove(filePath)
			return
		}
		writer.Write(arr[:len])
		hash.Write(arr[:len])
		writer.Flush()
	}
	f.Close()

	chunk = model.Chunk{
		MD5:        hex.EncodeToString(hash.Sum(nil)),
		FileName:   req.FileName,
		ChunkIndex: req.ChunkIndex,
		Size:       Handler.Size,
		Ext:        req.Ext,
		FilePath:   fmt.Sprintf("%d/%s_%d", id, req.FileName, req.ChunkIndex),
		AuthorID:   uint(id),
	}
	chunk.ID = 0
	err = l.svcCtx.DB.Create(&chunk).Error
	resp.Message = "上传成功"
	return
}
