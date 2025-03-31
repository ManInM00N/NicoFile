package comment

import (
	"context"
	"gorm.io/gorm"
	config2 "main/config"
	"main/model"
	"main/nicofile/internal/svc"
	"main/nicofile/internal/types"
	"main/pkg/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentListLogic {
	return &CommentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentListLogic) CommentList(req *types.CommentListRequest) (resp *types.CommentListResponse, err error) {
	resp = &types.CommentListResponse{
		Error:   false,
		List:    make([]types.Comment, 0),
		Message: "",
	}
	tot := int64(0)
	l.svcCtx.DB.Model(&model.Comment{}).Where("article_id = ?", req.ArticleId).Count(&tot)
	pages := (int(tot) + config2.PageSize - 1) / config2.PageSize
	req.Page = min(req.Page, pages)
	var list []model.Comment
	err2 := l.svcCtx.DB.Model(&model.Comment{}).
		Where("article_id = ? and parent_id IS NULL", req.ArticleId).
		Preload("Author").
		Preload("Replies", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Author").Order("created_at ASC")
		}).
		Find(&list).Error
	if err2 != nil {
		util.Log.Errorf("query list error: %v\n", err2)
	}
	for _, i := range list {
		tmp := types.Comment{
			Id:        int64(i.ID),
			CreatedAt: i.CreatedAt.Format("2006-01-02 15:04:05"),
			IP:        i.IP,
			Content:   i.Content,
			ArticleId: int64(i.ArticleID),
			AuthorId:  int64(i.AuthorId),
			Author:    i.Author.Username,
			Comments:  make([]types.Comment, 0),
		}
		for _, j := range i.Replies {
			tmp.Comments = append(tmp.Comments, types.Comment{
				Id:        int64(j.ID),
				CreatedAt: j.CreatedAt.Format("2006-01-02 15:04:05"),
				IP:        j.IP,
				Content:   j.Content,
				ArticleId: int64(i.ArticleID),
				AuthorId:  int64(j.AuthorId),
				Author:    j.Author.Username,
				ParentId:  int64(i.ID),
			})
		}
		resp.List = append(resp.List, tmp)
	}
	resp.Num = len(resp.List)
	resp.AllPages = pages
	resp.Page = req.Page
	return
}
