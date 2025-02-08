package file

import (
	"context"
	"encoding/json"
	"fmt"
	"main/model"
	"os"
	"path/filepath"

	"main/nicofile/internal/svc"
	"main/nicofile/internal/types"

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
	l.svcCtx.DB.Model(&model.File{}).Where("id = ?", req.FileId).First(&file)
	id, _ := l.ctx.Value("UserId").(json.Number).Int64()
	fmt.Println(id, file.ID)
	if (file.AuthorID != uint(id) && id <= 0) || file.ID == 0 {
		resp.Message = "无权删除"
		resp.Error = true
		return
	}
	if err2 := l.svcCtx.DB.Unscoped().Delete(&file).Error; err2 != nil {
		resp.Message = "删除失败或者文件不存在"
		resp.Error = true
		fmt.Println(file.ID)
	} else {
		os.Remove(filepath.Join(l.svcCtx.Config.StoragePath, file.FilePath))
	}
	return
}
