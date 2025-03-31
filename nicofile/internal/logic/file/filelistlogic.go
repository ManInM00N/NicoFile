package file

import (
	"context"
	"gorm.io/gorm"
	config2 "main/config"
	"main/model"
	"main/pkg/util"
	"strconv"

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
	l.svcCtx.DB.Model(&model.File{}).Count(&tot)
	pages := (int(tot) + config2.PageSize - 1) / config2.PageSize
	req.Page = min(req.Page, pages)
	offset := (req.Page - 1) * config2.PageSize
	var list []model.File
	//id, _ := l.ctx.Value("UserId").(json.Number).Int64()
	err2 := l.svcCtx.DB.Model(&model.File{}).Preload("Author", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, username").Where(" ?", true)
	}).
		Select("id,file_name,file_path,size,author_id,md5,ext,description,download_times,created_at").
		Where("  ?", true).
		Offset(offset).
		Limit(config2.PageSize).
		Find(&list).Error
	if err2 != nil {
		util.Log.Errorf("query list error: %v\n", err2)
	}
	for _, i := range list {
		tmp := types.File{
			Id:             i.ID,
			Name:           i.FileName,
			Path:           i.FilePath,
			Size:           i.Size,
			PosterId:       i.AuthorID,
			PosterName:     i.Author.Username,
			MD5:            i.MD5,
			Ext:            i.Ext,
			Desc:           i.Description,
			DonwloadCounts: i.DownloadTimes,
			CreatedAt:      i.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if !l.svcCtx.Config.Redis.Disabled {
			result, _ := l.svcCtx.Rdb.HGet(context.Background(), "file:"+strconv.Itoa(int(i.ID)), "download_times").Result()
			v, _ := strconv.Atoi(result)
			tmp.DonwloadCounts = int64(v)
		}
		resp.List = append(resp.List, tmp)
	}

	resp.Num = len(resp.List)
	resp.AllPages = pages
	resp.Page = req.Page
	return
}
