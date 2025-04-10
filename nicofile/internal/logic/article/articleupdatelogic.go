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

type ArticleUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewArticleUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleUpdateLogic {
	return &ArticleUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ArticleUpdateLogic) ArticleUpdate(req *types.ArticleUpdateRequest) (resp *types.ArticleResponse, err error) {
	resp = &types.ArticleResponse{
		Error: false,
	}
	id := l.ctx.Value("UserId").(json.Number).String()
	res := l.svcCtx.Rdb.HGet(context.Background(), fmt.Sprintf("article:%d", req.Id), "AuId").Val()
	if res != id && id != "1" {
		resp.Error = true
		resp.Message = "无权修改"
		return
	}
	//l.svcCtx.Rdb.HExpire(context.Background(), fmt.Sprintf("article:%d", req.Id), 60*60)
	if err2 := l.svcCtx.DB.Model(&model.Article{}).Where("id = ?", req.Id).Update("title", req.Title).Update("content", req.Content).Error; err2 != nil {
		resp.Error = true
		resp.Message = "更新失败"
		return
	}
	l.svcCtx.Rdb.HSet(context.Background(), fmt.Sprintf("article:%d", req.Id), "title", req.Title, "content", req.Content).Err()
	return
}
