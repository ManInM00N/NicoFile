package file

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"main/model"

	"main/nicofile/internal/svc"
	"main/nicofile/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
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
	id, _ := l.ctx.Value("UserId").(json.Number).Int64()
	if l.svcCtx.DB.Model(&model.File{}).Where("md5 = ? and file_name = ? and is_chunk = false and author_id = ?", req.FileMd5, req.FileName+req.Ext, id).Count(&num); num >= 1 {
		resp.Error = false
		resp.Message = "文件已存在"
		resp.Accept = req.ChunkNum
		return
	}

	for i, chunk := range req.MD5 {
		err = l.svcCtx.DB.Model(&model.File{}).Where("md5 = ? and file_name = ? and is_chunk = true and author_id = ?", chunk, fmt.Sprintf("%s_%d", req.FileName, i), id).Count(&num).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			resp.Error = true
			return
		}
		if num == 0 {
			resp.Accept = i
			return
		}
	}
	return
}
