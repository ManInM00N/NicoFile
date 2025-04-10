package file

import (
	"context"
	"encoding/json"
	"main/model"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"main/nicofile/internal/svc"
	"main/nicofile/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	chunkPool = sync.Pool{
		New: func() interface{} {
			return make([]model.Chunk, 0)
		},
	}
)

type CheckChunkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckChunkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckChunkLogic {
	return &CheckChunkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckChunkLogic) CheckChunk(req *types.CheckChunkRequest) (resp *types.CheckChunkResponse, err error) {
	resp = &types.CheckChunkResponse{
		Error:  false,
		Accept: req.ChunkNum,
	}
	var num int64
	idv := l.ctx.Value("UserId").(json.Number)
	id, _ := idv.Int64()
	if l.svcCtx.DB.Model(&model.File{}).Where("md5 = ? and file_name = ? and  author_id = ? and ext = ?", req.FileMd5, req.FileName+req.Ext, id, req.Ext).Count(&num); num >= 1 {
		resp.Error = false
		resp.Message = "文件已存在"
		resp.Accept = req.ChunkNum
		return
	}
	indexArr := make([]int, req.ChunkNum)
	for i := range req.ChunkNum {
		indexArr[i] = i
	}
	chunks := chunkPool.Get().([]model.Chunk)
	l.svcCtx.DB.Model(&model.Chunk{}).Select("id,chunk_index,md5,file_path").Where("file_name = ? and author_id = ? and ext = ? ", req.FileName, id, req.Ext).Find(&chunks)

	sort.Slice(chunks, func(i, j int) bool {
		return chunks[i].ChunkIndex < chunks[j].ChunkIndex
	})
	for i, chunk := range chunks {
		if chunk.ChunkIndex >= req.ChunkNum {
			break
		}
		if chunk.MD5 != req.MD5[i] || chunk.ChunkIndex != i {
			resp.Accept = i
			break
		}
	}
	if resp.Accept == req.ChunkNum && len(chunks) < req.ChunkNum {
		resp.Accept = len(chunks)
	}
	for _, chunk := range chunks {
		if chunk.ChunkIndex >= req.ChunkNum {
			os.Remove(filepath.Join(l.svcCtx.Config.ChunkStorePath, chunk.FilePath))
		}
	}
	if len(chunks) == 0 {
		resp.Accept = 0
	}
	l.svcCtx.DB.Unscoped().Model(&model.Chunk{}).Where("file_name = ? and author_id = ? and ext = ? and chunk_index >= ?", req.FileName, id, req.Ext, req.ChunkNum).Delete(&model.Chunk{})
	return
}
