package file

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"main/model"
	"main/nicofile/internal/svc"
	"main/nicofile/internal/types"
	"os"
	"path/filepath"

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
	resp = &types.MergeChunkResponse{Error: false}
	id := l.ctx.Value("UserId").(json.Number)
	path := filepath.Join(req.FileName + "_" + id.String() + req.MD5 + req.Ext)
	file, _ := os.OpenFile(filepath.Join(l.svcCtx.Config.StoragePath, path), os.O_CREATE|os.O_WRONLY, 0666)
	defer file.Close()
	id_v, _ := id.Int64()
	writer := bufio.NewWriter(file)
	for i := 0; i < req.ChunkNum; i++ {
		f, _ := os.OpenFile(fmt.Sprintf("%s/%s_%d", l.svcCtx.Config.ChunkStorePath, req.FileName, i), os.O_CREATE|os.O_RDONLY, 0666)
		reader := bufio.NewReader(f)
		_, err = io.Copy(writer, reader)
		if err != nil {
			resp.Error = true
			resp.Message = "合并分片文件失败"
			os.Remove(path)
			return resp, nil
		}
		writer.Flush()
		f.Close()
	}
	for i := 0; i < req.ChunkNum; i++ {
		l.svcCtx.DB.Unscoped().Delete(&model.File{}, "file_name = ? and is_chunk = true", fmt.Sprintf("%s_%d", req.FileName, i))
		os.Remove(fmt.Sprintf("%s/%s_%d", l.svcCtx.Config.ChunkStorePath, req.FileName, i))
	}
	l.svcCtx.DB.Create(&model.File{
		MD5:         req.MD5,
		FileName:    req.FileName + req.Ext,
		IsChunk:     false,
		Size:        req.Size,
		Ext:         req.Ext,
		FilePath:    path,
		AuthorID:    uint(id_v),
		Description: req.Description,
	})
	resp.Message = "合并分片文件成功"
	return
}
