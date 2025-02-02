package file

import (
	"context"
	"fmt"
	config2 "main/config"
	"main/model"

	"main/nicofile/internal/svc"
	"main/nicofile/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileListLogic {
	return &FileListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileListLogic) FileList(req *types.FileListRequest) (resp *types.FileListResponse, err error) {
	resp = &types.FileListResponse{
		List:  make([]types.File, 0),
		Error: false,
	}
	tot := int64(0)
	l.svcCtx.DB.Model(&types.File{}).Where("is_chunk = 0").Count(&tot)
	pages := (int(tot) + config2.PageSize - 1) / config2.PageSize
	req.Page = min(req.Page, pages)
	offset := (req.Page - 1) * config2.PageSize
	var list []model.File
	l.svcCtx.DB.Model(&types.File{}).Preload("Author").Where("is_chunk = 0").Offset(offset).Limit(config2.PageSize).Find(&list)
	for _, i := range list {
		resp.List = append(resp.List, types.File{
			Id:         i.ID,
			Name:       i.FileName,
			Path:       i.FilePath,
			Size:       i.Size,
			PosterId:   i.AuthorID,
			PosterName: i.Author.Username,
			MD5:        i.MD5,
			Ext:        i.Ext,
			Desc:       i.Description,
			CreatedAt:  i.CreatedAt.Format("2006-01-02 15:04:05"),
		})
		fmt.Println(i.Author)
	}
	resp.Num = len(resp.List)
	return
}
