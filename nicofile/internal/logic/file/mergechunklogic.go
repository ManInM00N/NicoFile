package file

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"google.golang.org/protobuf/proto"
	"io"
	"main/model"
	"main/nicofile/internal/svc"
	"main/nicofile/internal/types"
	"main/pkg/util"
	"main/server/proto/kafka"
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
	fileName := filepath.Join(req.FileName + "_" + id.String() + req.MD5 + req.Ext)
	fileStorePath := filepath.Join(l.svcCtx.Config.StoragePath, id.String())
	id_v, _ := id.Int64()
	os.MkdirAll(fileStorePath, os.ModePerm)
	chunkStorePath := fmt.Sprintf("%s/%s", l.svcCtx.Config.ChunkStorePath, id.String())
	file, errt := os.OpenFile(filepath.Join(fileStorePath, fileName), os.O_CREATE|os.O_WRONLY, 0666)
	if errt != nil {
		util.Log.Errorf("create file error: %v", errt)
	}
	writer := bufio.NewWriter(file)

	indexArr := make([]int, req.ChunkNum)
	broken := false
	for i := 0; i < req.ChunkNum; i++ {
		indexArr[i] = i
		f, err2 := os.OpenFile(fmt.Sprintf("%s/%s_%d", chunkStorePath, req.FileName, i), os.O_RDONLY, 0666)
		if err2 != nil {
			if errors.Is(err2, os.ErrNotExist) {
				util.Log.Errorf("file open failed %v: %v", req.FileName, err2)
			}
			broken = true
			break
		}
		reader := bufio.NewReader(f)
		_, err2 = io.Copy(writer, reader)
		if err2 != nil {
			resp.Error = true
			resp.Message = "合并分片文件失败"
			os.Remove(filepath.Join(fileStorePath, fileName))
			return resp, nil
		}
		writer.Flush()
		f.Close()
	}
	file.Close()

	if broken {
		resp.Error = true
		resp.Message = "合并分片文件失败"
		os.Remove(filepath.Join(fileStorePath, fileName))
		return
	}
	num := l.svcCtx.DB.Unscoped().Delete(&model.Chunk{}, "file_name = ? and ext = ? and author_id = ?", req.FileName, req.Ext, id_v).RowsAffected
	if num != int64(req.ChunkNum) {

	}
	for i := 0; i < req.ChunkNum; i++ {
		os.Remove(fmt.Sprintf("%s/%s_%d", chunkStorePath, req.FileName, i))
	}
	err2 := l.svcCtx.DB.Create(&model.File{
		MD5:         req.MD5,
		FileName:    req.FileName + req.Ext,
		IsChunk:     false,
		Size:        req.Size,
		Ext:         req.Ext,
		FilePath:    filepath.Join(id.String(), fileName),
		AuthorID:    uint(id_v),
		Description: req.Description,
	}).Error
	if err2 != nil {
		resp.Error = true
		resp.Message = "合并分片文件失败"
		os.Remove(filepath.Join(fileStorePath, fileName))
		util.Log.Errorf("create file error: %v", err2)
		return
	} else if !l.svcCtx.Config.Redis.Disabled {
		key := fmt.Sprintf("file:%d", id_v)
		l.svcCtx.Rdb.HSet(context.Background(), key, "upload_times", 0, "AuId", id_v, "description", req.Description)
	}

	if !l.svcCtx.Config.Kafka.Disabled {
		event := &kafka.FileMonitor{
			Message: "A file has been uploaded",
			Warning: false,
			UserId:  uint32(id_v),
		}
		data, _ := proto.Marshal(event)
		(*l.svcCtx.Producer).Input() <- &sarama.ProducerMessage{
			Topic: "file-monitor",
			Value: sarama.ByteEncoder(data),
		}
	}
	resp.Message = "合并分片文件成功"

	return
}
