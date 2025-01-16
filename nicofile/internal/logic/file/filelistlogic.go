package file

import (
	"context"
	config2 "main/config"

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
	l.svcCtx.DB.Model(&types.File{}).Count(&tot)
	pages := (int(tot) + config2.PageSize - 1) / config2.PageSize
	req.Page = min(req.Page, pages)
	offset := (req.Page - 1) * config2.PageSize
	l.svcCtx.DB.Model(&types.File{}).Offset(offset).Limit(config2.PageSize).Find(&resp.List)
	resp.Num = len(resp.List)
	return
}
