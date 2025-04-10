package handler

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	config2 "main/config"
	"main/model"
	"main/pkg/util"
	CacheRedis "main/redis"
	"main/server/proto/articleRank"
	"strconv"
	"sync"
	"time"
)

type ArticleRankServiceServer struct {
	articleRank.UnimplementedArticleRankServiceServer
	currentWindow string
	windowMutex   sync.Mutex
}

func (s *ArticleRankServiceServer) updateCurrentWindow() {
	s.windowMutex.Lock()
	defer s.windowMutex.Unlock()
	rdb := CacheRedis.GetRdb()
	now := time.Now()
	windowStart := now.Truncate(10 * time.Minute)
	s.currentWindow = fmt.Sprintf("article:hotness:window:%d", windowStart.Unix())

	// 设置当前窗口key的过期时间为20分钟(比窗口时间长)
	rdb.Expire(context.Background(), "article:hotness:current_window", 20*time.Minute)
}
func (s *ArticleRankServiceServer) RecordView(articleID uint) error {
	// 1. 更新MySQL总访问量
	rdb := CacheRedis.GetRdb()
	// 2. 记录到当前10分钟窗口(访问权重=1)
	if err := rdb.ZIncrBy(context.Background(), "article:hotness:current_window", 1,
		fmt.Sprintf("%d", articleID)).Err(); err != nil {
		return err
	}

	// 3. 同时更新总排行榜(访问权重=1)
	if err := rdb.ZIncrBy(context.Background(), "article:hotness:leaderboard", 1,
		fmt.Sprintf("%d", articleID)).Err(); err != nil {
		return err
	}

	return nil
}

func (s *ArticleRankServiceServer) RecordLike(articleID uint) error {
	// 1. 更新MySQL总点赞量
	rdb := CacheRedis.GetRdb()
	// 2. 记录到当前10分钟窗口(点赞权重=3)
	if err := rdb.ZIncrBy(context.Background(), "article:hotness:current_window", 3,
		fmt.Sprintf("%d", articleID)).Err(); err != nil {
		return err
	}

	// 3. 同时更新总排行榜(点赞权重=3)
	if err := rdb.ZIncrBy(context.Background(), "article:hotness:leaderboard", 3,
		fmt.Sprintf("%d", articleID)).Err(); err != nil {
		return err
	}

	return nil
}
func (s *ArticleRankServiceServer) startWindowRotation() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for {
		<-ticker.C
		s.rotateWindow()
	}
}

func (s *ArticleRankServiceServer) rotateWindow() {
	// 1. 获取当前窗口数据
	rdb := CacheRedis.GetRdb()
	currentWindow := s.currentWindow
	windowData, err := rdb.ZRangeWithScores(context.Background(),
		"article:hotness:current_window", 0, -1).Result()
	if err != nil {
		util.Log.Printf("Failed to get current window data: %v", err)
		return
	}

	// 2. 重命名当前窗口为历史窗口
	if len(windowData) > 0 {
		// 使用MULTI保证原子性
		pipe := rdb.TxPipeline()

		// 重命名当前窗口
		pipe.Rename(context.Background(), "article:hotness:current_window", currentWindow)

		// 添加到窗口索引
		windowTime := time.Now().Truncate(10 * time.Minute).Add(-10 * time.Minute)
		pipe.ZAdd(context.Background(), "article:hotness:windows", redis.Z{
			Score:  float64(windowTime.Unix()),
			Member: currentWindow,
		})

		// 设置历史窗口过期时间为15天
		pipe.Expire(context.Background(), currentWindow, 15*24*time.Hour)

		if _, err := pipe.Exec(context.Background()); err != nil {
			log.Printf("Failed to rotate window: %v", err)
		}
	}

	// 3. 清理过期窗口(15天前)
	fifteenDaysAgo := time.Now().AddDate(0, 0, -15).Unix()
	expiredKeys, err := rdb.ZRangeByScore(context.Background(),
		"article:hotness:windows", &redis.ZRangeBy{
			Min: "0",
			Max: strconv.FormatInt(fifteenDaysAgo, 10),
		}).Result()
	if err == nil && len(expiredKeys) > 0 {
		// 从索引中移除
		rdb.ZRemRangeByScore(context.Background(), "article:hotness:windows",
			"0", strconv.FormatInt(fifteenDaysAgo, 10))

		// 删除实际数据
		rdb.Del(context.Background(), expiredKeys...)

		// 重新计算总排行榜(因为移除了部分数据)
		s.recalculateLeaderboard()
	}

	// 4. 更新当前窗口
	s.updateCurrentWindow()
}

func (s *ArticleRankServiceServer) recalculateLeaderboard() {
	// 1. 获取所有有效窗口
	rdb := CacheRedis.GetRdb()
	fifteenDaysAgo := time.Now().AddDate(0, 0, -15).Unix()
	windows, err := rdb.ZRangeByScore(context.Background(),
		"article:hotness:windows", &redis.ZRangeBy{
			Min: strconv.FormatInt(fifteenDaysAgo, 10),
			Max: "+inf",
		}).Result()
	if err != nil || len(windows) == 0 {
		return
	}

	// 2. 使用Lua脚本合并所有窗口数据
	script := `
        -- 删除旧排行榜
        redis.call('DEL', KEYS[1])
        
        -- 合并所有窗口数据
        for i=2,#KEYS do
            local members = redis.call('ZRANGE', KEYS[i], 0, -1, 'WITHSCORES')
            for j=1,#members,2 do
                redis.call('ZINCRBY', KEYS[1], members[j+1], members[j])
            end
        end
    `

	// 准备参数: [leaderboardKey, window1, window2, ...]
	args := make([]string, len(windows)+1)
	args[0] = "article:hotness:leaderboard"
	copy(args[1:], windows)

	if _, err := rdb.Eval(context.Background(), script, args).Result(); err != nil {
		log.Printf("Failed to recalculate leaderboard: %v", err)
	}
}

func (s *ArticleRankServiceServer) GetArticleRank(ctx context.Context, req *articleRank.GetArticleRankRequest) (*articleRank.GetArticleRankResponse, error) {
	rdb := CacheRedis.GetRdb()
	articleIDs, err := rdb.ZRevRange(context.Background(),
		"article:hotness:leaderboard", 0, int64(req.ArticleNum-1)).Result()
	if err != nil {
		return nil, err
	}

	// 2. 从MySQL查询文章详情
	var articles []*model.Article
	db := config2.GetDB()
	if err2 := db.Where("id IN ?", articleIDs).Find(&articles).Error; err2 != nil {
		return nil, err2
	}
	fmt.Println(articles)
	// 3. 按照热度排序
	idToArticle := make(map[string]*model.Article)
	for _, article := range articles {

		idToArticle[fmt.Sprintf("%d", article.ID)] = article
	}

	result := make([]*articleRank.Article, 0, len(articleIDs))
	for _, id := range articleIDs {
		if article, ok := idToArticle[id]; ok {
			result = append(result, &articleRank.Article{
				ArticleId:    uint32(article.ID),
				ArticleTitle: article.Title,
			})
		}
	}
	res := &articleRank.GetArticleRankResponse{
		Success:     true,
		ArticleList: result,
	}
	return res, nil
}

func NewArticleRankService() *ArticleRankServiceServer {
	// 初始化当前窗口
	s := &ArticleRankServiceServer{}
	s.updateCurrentWindow()

	// 启动定时任务
	go s.startWindowRotation()
	return s
}
