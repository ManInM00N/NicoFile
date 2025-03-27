package article

import (
	"context"
	config2 "main/config"
	"main/model"
	"main/pkg/util"
	"strconv"

	"main/nicofile/internal/svc"
	"main/nicofile/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewArticleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleListLogic {
	return &ArticleListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ArticleListLogic) ArticleList(req *types.ArticleListRequest) (resp *types.ArticleListResponse, err error) {
	resp = &types.ArticleListResponse{
		Error:   false,
		List:    make([]types.Article, 0),
		Message: "",
	}
	tot := int64(0)
	l.svcCtx.DB.Model(&model.Article{}).Count(&tot)
	pages := (int(tot) + config2.PageSize - 1) / config2.PageSize
	req.Page = min(req.Page, pages)
	offset := (req.Page - 1) * config2.PageSize
	var list []model.Article
	err2 := l.svcCtx.DB.Model(&model.Article{}).Preload("Author").
		Offset(offset).
		Limit(config2.PageSize).
		Find(&list).Error
	if err2 != nil {
		util.Log.Errorf("query list error: %v\n", err2)
	}
	for _, i := range list {
		tmp := types.Article{
			Id:         int64(i.ID),
			CreatedAt:  i.CreatedAt.Format("2006-01-02 15:04:05"),
			Content:    i.Content,
			Title:      i.Title,
			View:       i.View,
			Like:       i.Like,
			AuthorId:   int64(i.AuthorID),
			AuthorName: i.Author.Username,
			Cover:      i.Cover,
		}
		result, _ := l.svcCtx.Rdb.HGet(context.Background(), "article:"+strconv.Itoa(int(i.ID)), "view").Result()
		v, _ := strconv.Atoi(result)
		tmp.View = int64(v)
		result, _ = l.svcCtx.Rdb.HGet(context.Background(), "article:"+strconv.Itoa(int(i.ID)), "like").Result()
		v, _ = strconv.Atoi(result)
		tmp.Like = int64(v)
		resp.List = append(resp.List, tmp)
	}
	resp.Num = len(resp.List)
	resp.AllPages = pages
	resp.Page = req.Page
	return
}
