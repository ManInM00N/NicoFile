package file

import (
	"bufio"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"main/model"
	"main/pkg/encrypt"
	"os"

	"main/nicofile/internal/svc"
	"main/nicofile/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
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

func (l *UploadChunkLogic) UploadChunk(req *types.UploadChunkRequest) (resp *types.UploadChunkResponse, err error) {
	// todo: add your logic here and delete this line

	num := int64(0)
	resp.Error = false
	l.svcCtx.DB.Where("md5 = ?", req.MD5).Count(&num)
	if num == 1 {
		resp.Message = "文件已存在"
		return
	}
	f, _ := os.OpenFile(fmt.Sprintf("%s/%s_%d", l.svcCtx.Config.ChunkStorePath, req.FileName, req.ChunkIndex), os.O_CREATE|os.O_WRONLY, 0666)
	writer := bufio.NewWriter(f)
	_, err = writer.Write(req.Chunk)

	if err != nil {
		logrus.Errorln(err)
		resp.Error = true
		resp.Message = "写入分片文件失败"
		return resp, nil
	}
	if err = writer.Flush(); err != nil {
		return nil, err
	}
	file := model.File{
		MD5:      encrypt.Md5Sum(req.Chunk),
		FileName: req.FileName,
		IsChunk:  true,
		Size:     int64(len(req.Chunk)),
	}
	l.svcCtx.DB.Create(&file)
	resp.Message = "上传成功"
	return
}
