package article

import (
	"context"
	"fmt"
	"main/nicofile/internal/svc"
	"main/nicofile/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewArticleDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleDetailLogic {
	return &ArticleDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ArticleDetailLogic) ArticleDetail(req *types.ArticleRequest) (resp *types.ArticleResponse, err error) {
	resp = &types.ArticleResponse{
		Error:   false,
		Message: "",
	}
	res, _ := l.svcCtx.Rdb.HGetAll(context.Background(), fmt.Sprintf("article:%d", req.Id)).Result()

	if res["AuId"] != "" {
		resp.Content = res["content"]
		resp.Title = res["title"]
		resp.CreatedAt = res["creat_at"]
		resp.View = res["view"]
		l.svcCtx.Rdb.HIncrBy(context.Background(), fmt.Sprintf("article:%d", req.Id), "view", 1)
	} else {
		resp.Error = true
		resp.Message = "文章不存在"
	}

	return
}
