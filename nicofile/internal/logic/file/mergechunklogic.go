package file

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"main/model"
	"os"

	"main/nicofile/internal/svc"
	"main/nicofile/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MergeChunkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMergeChunkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MergeChunkLogic {
	return &MergeChunkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MergeChunkLogic) MergeChunk(req *types.MergeChunkRequest) (resp *types.MergeChunkResponse, err error) {
	// todo: add your logic here and delete this line
	file, _ := os.OpenFile(req.FileName, os.O_CREATE|os.O_WRONLY, 0666)
	writer := bufio.NewWriter(file)
	for i := 0; i < req.ChunkNum; i++ {
		f, _ := os.OpenFile(fmt.Sprintf("%s/%s_%d", l.svcCtx.Config.ChunkStorePath, req.FileName, i), os.O_CREATE|os.O_RDONLY, 0666)
		reader := bufio.NewReader(f)
		_, err = io.Copy(writer, reader)
		if err != nil {
			resp.Error = true
			resp.Message = "合并分片文件失败"
			file.Close()
			os.Remove(req.FileName)
			return resp, nil
		}
		writer.Flush()
		f.Close()
	}
	file.Close()
	for i := 0; i < req.ChunkNum; i++ {
		l.svcCtx.DB.Delete(&model.File{}, "file_name = ? and is_chunk = true", fmt.Sprintf("%s/%s_%d", l.svcCtx.Config.ChunkStorePath, req.FileName, i))
	}
	l.svcCtx.DB.Create(&model.File{
		MD5:      req.MD5,
		FileName: req.FileName,
		IsChunk:  false,
		Size:     req.Size,
		Ext:      req.Ext,
	})
	resp.Message = "合并分片文件成功"
	return
}
