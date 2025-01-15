package file

import (
	"bufio"
	"context"
	"crypto/md5"
	"encoding/hex"
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
	l.svcCtx.DB.Model(&model.File{}).Where("md5 = ? and file_name = ?", req.MD5, fmt.Sprintf("%s_%d", req.FileName, req.ChunkIndex)).Count(&num)
	if num >= 1 {
		resp.Message = "文件已存在"
		return
	}
	os.MkdirAll(l.svcCtx.Config.ChunkStorePath, os.ModePerm)
	f, err := os.OpenFile(fmt.Sprintf("%s/%s_%d", l.svcCtx.Config.ChunkStorePath, req.FileName, req.ChunkIndex), os.O_CREATE|os.O_WRONLY, 0666)
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

	for {
		len, err := (*File).Read(arr)
		if err == io.EOF || len == 0 {
			break
		}
		writer.Write(arr[:len])
		hash.Write(arr[:len])
		writer.Flush()
	}
	(*File).Close()
	f.Close()
	file := model.File{
		MD5:      hex.EncodeToString(hash.Sum(nil)),
		FileName: fmt.Sprintf("%s_%d", req.FileName, req.ChunkIndex),
		IsChunk:  true,
		Size:     Handler.Size,
	}
	l.svcCtx.DB.Create(&file)
	resp.Message = "上传成功"
	return
}
