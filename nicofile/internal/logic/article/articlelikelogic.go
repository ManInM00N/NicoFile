package article

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"main/model"
	"strconv"

	"main/nicofile/internal/svc"
	"main/nicofile/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleLikeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewArticleLikeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleLikeLogic {
	return &ArticleLikeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ArticleLikeLogic) ArticleLike(req *types.ArticleLikeRequest) (resp *types.Response, err error) {
	resp = &types.Response{}
	res, _ := l.svcCtx.Rdb.HGet(context.Background(), fmt.Sprintf("article:%d", req.Id), "like").Result()
	if res == "" {
		var art model.Article
		if errors.Is(l.svcCtx.DB.Model(&model.Article{}).Where("id = ?", req.Id).First(&art).Error, gorm.ErrRecordNotFound) {
			resp.Message = "文章不存在"
			return
		}
		l.svcCtx.Rdb.HSet(context.Background(), fmt.Sprintf("article:%d", req.Id), "AuId", art.AuthorID, "content", art.Content, "title", art.Title, "creat_at", art.CreatedAt.Format("2006-01-02 15:04:05"), "view", art.View, "like", art.Like+1)
		l.svcCtx.Rdb.ZIncrBy(context.Background(), "article:hotness:current_window", 1, fmt.Sprintf("%d", req.Id))
		l.svcCtx.Rdb.ZIncrBy(context.Background(), "article:hotness:leaderboard", 1, strconv.FormatInt(req.Id, 10))
		resp.Message = "点赞成功"
	} else {
		l.svcCtx.Rdb.HIncrBy(context.Background(), fmt.Sprintf("article:%d", req.Id), "like", 1)
		l.svcCtx.Rdb.ZIncrBy(context.Background(), "article:hotness:current_window", 1, fmt.Sprintf("%d", req.Id))
		l.svcCtx.Rdb.ZIncrBy(context.Background(), "article:hotness:leaderboard", 1, strconv.FormatInt(req.Id, 10))
		resp.Message = "点赞成功"
	}
	return
}
