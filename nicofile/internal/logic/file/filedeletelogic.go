package file

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
	"main/model"
	"main/nicofile/internal/svc"
	"main/nicofile/internal/types"
	"main/pkg/util"
	"main/server/proto/kafka"
	"os"
	"path/filepath"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileDeleteLogic {
	return &FileDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileDeleteLogic) FileDelete(req *types.FileDeleteRequest) (resp *types.FileDeleteResponse, err error) {
	resp = &types.FileDeleteResponse{
		Error: false,
	}
	var file model.File
	ctx := context.Background()
	id, _ := l.ctx.Value("UserId").(json.Number).Int64()

	if errt := l.svcCtx.DB.Model(&model.File{}).Select("id,author_id,download_times,description,file_path").Where("id = ?", req.FileId).First(&file).Error; errt != nil {
		if errors.Is(errt, gorm.ErrRecordNotFound) {
			resp.Message = "文件不存在"
			resp.Error = true
			return
		} else if errt != nil {
			resp.Message = "数据库错误"
			resp.Error = true
			util.Log.Errorf("Failed to query file %d: %v", req.FileId, errt)
			return
		}
	}
	if (file.AuthorID != uint(id) && id > 0) || file.ID == 0 {
		resp.Message = "无权删除"
		resp.Error = true
		return
	}
	if !l.svcCtx.Config.Redis.Disabled {
		key := fmt.Sprintf("file:%d", req.FileId)
		err2 := l.svcCtx.Rdb.Del(ctx, key).Err()
		if err2 != nil {
			resp.Message = "删除失败"
			resp.Error = true
			util.Log.Errorf("Failed to delete file %d from Redis: %v", file.ID, err2)
			return
		}
	}
	if err2 := l.svcCtx.DB.Unscoped().Delete(&file).Error; err2 != nil {
		resp.Message = "删除失败或者文件不存在"
		resp.Error = true
	} else {
		if err2 := os.Remove(filepath.Join(l.svcCtx.Config.StoragePath, file.FilePath)); err2 != nil {
			resp.Message = "删除失败或者文件不存在"
			resp.Error = true
			util.Log.Errorln("Failed to delete file %d from disk, %v", file.ID, err2)
		} else if !l.svcCtx.Config.Kafka.Disabled {
			var event = &kafka.FileMonitor{
				Message: "A file has been deleted",
				Warning: false,
				UserId:  uint32(id),
			}
			data, _ := proto.Marshal(event)
			(*l.svcCtx.Producer).Input() <- &sarama.ProducerMessage{
				Topic: "file-monitor",
				Value: sarama.ByteEncoder(data),
			}
		}
	}

	return
}
