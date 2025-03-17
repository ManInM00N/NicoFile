package article

import (
	"context"
	"encoding/json"
	"fmt"
	"main/model"

	"main/nicofile/internal/svc"
	"main/nicofile/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewArticleCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleCreateLogic {
	return &ArticleCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ArticleCreateLogic) ArticleCreate(req *types.ArticleCreateRequest) (resp *types.ArticleResponse, err error) {
	resp = &types.ArticleResponse{
		Error:   false,
		Message: "",
	}
	idv := l.ctx.Value("userId").(json.Number)
	id, _ := idv.Int64()
	Art := model.Article{
		Title:    req.Title,
		Content:  req.Content,
		AuthorID: uint(id),
	}
	err2 := l.svcCtx.DB.Create(&Art).Error
	if err2 != nil {
		resp.Error = true
		resp.Message = "创建失败error creating"
		return
	}
	l.svcCtx.Rdb.HSet(context.Background(), fmt.Sprintf("article:%d", Art.ID), "title", Art.Title, "AuId", id, "content", Art.Content, "creat_at", Art.CreatedAt).Err()

	return
}
