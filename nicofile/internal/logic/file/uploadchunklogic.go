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
	// todo: add your logic here and delete this line
	num := int64(0)
	resp = &types.UploadChunkResponse{
		Error: false,
	}
	id, _ := l.ctx.Value("UserId").(json.Number).Int64()
	l.svcCtx.DB.Model(&model.File{}).Preload("Author").Where("md5 = ? and file_name = ? and author_id = ?", req.MD5, fmt.Sprintf("%s_%d", req.FileName, req.ChunkIndex), id).Count(&num)
	if num >= 1 {
		resp.Message = "文件已存在"
		return
	}
	os.MkdirAll(l.svcCtx.Config.ChunkStorePath, os.ModePerm)
	storagePath := fmt.Sprintf("%s/%s_%d", l.svcCtx.Config.ChunkStorePath, req.FileName, req.ChunkIndex)
	f, err := os.OpenFile(storagePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		resp.Error = true
		resp.Message = "创建分片文件失败"
		return resp, nil
	}
	//bucket := ratelimit.NewBucket(20*time.Microsecond, 5)
	writer := bufio.NewWriter(f)
	//W := ratelimit.Writer(writer, bucket)

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
			os.Remove(storagePath)
			return
		}
		writer.Write(arr[:len])
		hash.Write(arr[:len])
		writer.Flush()
	}
	f.Close()

	file := model.File{
		MD5:      hex.EncodeToString(hash.Sum(nil)),
		FileName: fmt.Sprintf("%s_%d", req.FileName, req.ChunkIndex),
		IsChunk:  true,
		Size:     Handler.Size,
		FilePath: fmt.Sprintf("%s_%d", req.FileName, req.ChunkIndex),
		AuthorID: uint(id),
	}
	l.svcCtx.DB.Create(&file)
	resp.Message = "上传成功"
	return
}
