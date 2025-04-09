package article

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"main/model"
	"main/nicofile/internal/svc"
	"main/nicofile/internal/types"
	"strconv"
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
		resp.Cover = res["cover"]
		resp.Content = res["content"]
		resp.Title = res["title"]
		resp.CreatedAt = res["creat_at"]
		resp.View = res["view"]
		resp.Cover = res["cover"]
		resp.AuthorId = res["AuId"]
		resp.Like = res["like"]
		res, _ = l.svcCtx.Rdb.HGetAll(context.Background(), fmt.Sprintf("user:%s", resp.AuthorId)).Result()
		resp.AuthorName = res["username"]
		resp.ArticleId = strconv.Itoa(int(req.Id))
		l.svcCtx.Rdb.HIncrBy(context.Background(), fmt.Sprintf("article:%d", req.Id), "view", 1)
		l.svcCtx.Rdb.ZIncrBy(context.Background(), "article:hotness:current_window", 1, fmt.Sprintf("%d", req.Id))
		l.svcCtx.Rdb.ZIncrBy(context.Background(), "article:hotness:leaderboard", 1, strconv.FormatInt(req.Id, 10))
		//l.svcCtx.Rdb.HExpire(context.Background(), fmt.Sprintf("article:%d", req.Id), 60*60)
	} else {
		var art model.Article
		if err2 := l.svcCtx.DB.Model(&model.Article{}).Where("id = ?", req.Id).Preload("Author").First(&art).Error; err2 != nil || art.AuthorID == 0 {
			resp.Error = true
			resp.Message = "文章不存在"
			return
		}
		l.svcCtx.Rdb.HSet(context.Background(), fmt.Sprintf("user:%d", art.Author.ID), "username", art.Author.Username, "priority", art.Author.Priority, "password", art.Author.Password)
		l.svcCtx.Rdb.HSet(context.Background(), fmt.Sprintf("article:%d", req.Id), "AuId", art.AuthorID, "content", art.Content, "title", art.Title, "creat_at", art.CreatedAt.Format("2006-01-02 15:04:05"), "view", art.View+1, "like", art.Like, "cover", art.Cover)
		//l.svcCtx.Rdb.HExpire(context.Background(), fmt.Sprintf("article:%d", req.Id), 60*60)
		resp.Content = art.Content
		resp.Title = art.Title
		resp.Cover = art.Cover
		resp.CreatedAt = art.CreatedAt.Format("2006-01-02 15:04:05")
		resp.View = strconv.FormatInt(art.View+1, 10)
		resp.Like = strconv.FormatInt(art.Like, 10)
		resp.AuthorId = strconv.FormatInt(int64(art.AuthorID), 10)
		resp.AuthorName = art.Author.Username
		resp.ArticleId = strconv.Itoa(int(req.Id))
		l.svcCtx.Rdb.ZIncrBy(context.Background(), "article:hotness:current_window", 3, fmt.Sprintf("%d", req.Id))
		l.svcCtx.Rdb.ZIncrBy(context.Background(), "article:hotness:leaderboard", 3, strconv.FormatInt(req.Id, 10))
	}

	return
}
