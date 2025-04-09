package article

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ManInM00N/go-tool/statics"
	"github.com/zeromicro/go-zero/core/logx"
	"main/model"
	"main/nicofile/internal/svc"
	"main/nicofile/internal/types"
	"strconv"
)

type ArticleSearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewArticleSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleSearchLogic {
	return &ArticleSearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type Result struct {
	Hits struct {
		Total struct {
			Value int64 `json:"value"`
		} `json:"total"`
		Hits []struct {
			Source model.Article `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
	Took int `json:"took"`
}

func (l *ArticleSearchLogic) ArticleSearch(req *types.ArticleSearchRequest) (resp *types.ArticleListResponse, err error) {
	resp = &types.ArticleListResponse{
		Error: false,
		List:  make([]types.Article, 0),
	}
	var buf bytes.Buffer

	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":    req.Keyword,
				"fields":   []string{"title", "content"},
				"operator": "or", // title 或 content 匹配即可
			},
		},
		"sort": []map[string]interface{}{
			{"_score": map[string]interface{}{"order": "desc"}},     // 按相关性
			{"created_at": map[string]interface{}{"order": "desc"}}, // 按时间
		},
	}
	if err2 := json.NewEncoder(&buf).Encode(searchQuery); err2 != nil {
		return nil, fmt.Errorf("编码查询失败: %w", err2)
	}
	res, err2 := l.svcCtx.ES.Search(
		l.svcCtx.ES.Search.WithFrom((req.Page-1)*int(req.Size)), // 分页起始
		l.svcCtx.ES.Search.WithSize(int(req.Size)),
		l.svcCtx.ES.Search.WithIndex("articles"),
		l.svcCtx.ES.Search.WithBody(&buf),
		l.svcCtx.ES.Search.WithPretty(),
	)
	if err2 != nil {
		return nil, fmt.Errorf("搜索请求失败: %w", err2)
	}
	defer res.Body.Close()
	var result Result
	//fmt.Println(res.String())
	if res.IsError() {
		return nil, fmt.Errorf("搜索返回错误: %s", res.String())
	}
	if err2 := json.NewDecoder(res.Body).Decode(&result); err2 != nil {
		return nil, fmt.Errorf("解析结果失败: %w", err2)
	}
	for _, hit := range result.Hits.Hits {
		fmt.Println(hit.Source.CreatedAt)
		resp.List = append(resp.List, types.Article{
			Id:         int64(hit.Source.ID),
			Cover:      hit.Source.Cover,
			Title:      hit.Source.Title,
			Content:    hit.Source.Content,
			AuthorId:   int64(hit.Source.AuthorID),
			AuthorName: l.svcCtx.Rdb.HGet(context.Background(), "user:"+strconv.Itoa(int(hit.Source.AuthorID)), "username").Val(),
			View:       statics.StringToInt64(l.svcCtx.Rdb.HGet(context.Background(), "article:"+strconv.Itoa(int(hit.Source.ID)), "view").Val()),
			Like:       statics.StringToInt64(l.svcCtx.Rdb.HGet(context.Background(), "article:"+strconv.Itoa(int(hit.Source.ID)), "like").Val()),
			CreatedAt:  hit.Source.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return
}
