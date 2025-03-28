package article

import (
	"context"
	"errors"
	"fmt"
	"main/server/proto/articleRank"

	"main/nicofile/internal/svc"
	"main/nicofile/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleRankLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewArticleRankLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleRankLogic {
	return &ArticleRankLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ArticleRankLogic) ArticleRank(req *types.Null) (resp *types.ArticleListResponse, err error) {
	resp = &types.ArticleListResponse{
		Error: false,
	}
	// 从连接池获取连接
	conn, err2 := l.svcCtx.HotArticlePool.Get()
	if err2 != nil {
		return nil, errors.New("service is not available")
	}
	defer l.svcCtx.HotArticlePool.Put(conn)

	client := articleRank.NewArticleRankServiceClient(conn)
	res, err2 := client.GetArticleRank(context.Background(), &articleRank.GetArticleRankRequest{
		ArticleNum: 7,
	})

	if err2 != nil {
		fmt.Println(err2)
		return nil, errors.New("failed to connect service")
	}
	resp.List = make([]types.Article, 0)

	for _, article := range res.ArticleList {
		resp.List = append(resp.List, types.Article{Id: int64(article.ArticleId), Title: article.ArticleTitle})
	}
	return
}
