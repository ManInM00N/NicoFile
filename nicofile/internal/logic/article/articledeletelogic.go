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

type ArticleDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewArticleDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleDeleteLogic {
	return &ArticleDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ArticleDeleteLogic) ArticleDelete(req *types.ArticleDeleteRequest) (resp *types.ArticleDeleteResponse, err error) {
	resp = &types.ArticleDeleteResponse{
		Error: false,
	}
	res, _ := l.svcCtx.Rdb.HGet(context.Background(), fmt.Sprintf("article:%d", req.Id), "AuId").Result()
	id := l.ctx.Value("userId").(json.Number).String()
	if res != id {
		resp.Error = true
		resp.Message = "无权删除"
		return
	}

	if err2 := l.svcCtx.DB.Model(&model.Article{}).Where("id = ?", req.Id).Delete(&model.Article{}).Error; err2 != nil {
		resp.Error = true
		resp.Message = "删除失败"
		return
	}
	resp.Message = "删除成功"
	return
}
