package CacheRedis

import (
	"context"
	"fmt"
	"github.com/ManInM00N/go-tool/statics"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"main/model"
	"main/pkg/util"
	"strconv"
	"sync"
	"time"
)

var (
	Rdb  *redis.Client
	keyp = sync.Pool{
		New: func() interface{} {
			return make([]string, 0)
		},
	}
)

func GetRdb() *redis.Client {
	return Rdb
}
func InitRedis(host string, port int, password string, DB int, disabled bool) *redis.Client {
	if disabled {
		return nil
	}
	Rdb = redis.NewClient(&redis.Options{
		Addr:     host + strconv.Itoa(port),
		Password: password,
		DB:       DB,
	})
	_, err := Rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	return Rdb

}
func Transport(rdb *redis.Client, DB *gorm.DB) {

	ctx := context.Background()
	util.Log.Println("Starting sync from Redis to SQLite...")
	var cursor uint64
	counts := 0
	for {
		keys := keyp.Get().([]string)
		var err error
		keys, cursor, err = rdb.Scan(ctx, cursor, "file:*", 10000).Result()
		if err != nil {
			util.Log.Errorf("Failed to scan keys from Redis: %v", err)
			return
		}
		var list []model.File
		ids := make([]int, 0)
		for _, key := range keys {
			id := 0
			fmt.Sscanf(key, "file:%d", &id)
			ids = append(ids, id)
		}
		DB.Model(&model.File{}).Where("id in ?", ids).Find(&list)
		for _, i := range list {
			data, err := rdb.HGetAll(ctx, strconv.Itoa(int(i.ID))).Result()
			if err != nil {
				util.Log.Errorf("Failed to get hash data for key %d: %v", i.ID, err)
			}
			i.Description = data["description"]
			i.DownloadTimes = statics.StringToInt64(data["download_times"])
		}
		DB.Model(&model.File{}).Save(&list)
		//rdb.Del(ctx, keys...)
		counts += len(keys)
		keyp.Put(keys)
		if cursor == 0 {
			break
		}
	}
	util.Log.Println("Sync Files from Redis to SQLite finished, total:", counts)
	counts = 0
	cursor = 0
	for {
		keys := keyp.Get().([]string)
		var err error
		keys, cursor, err = rdb.Scan(ctx, cursor, "article:*", 10000).Result()
		if err != nil {
			util.Log.Errorf("Failed to scan keys from Redis: %v", err)
			return
		}
		var list []model.Article
		ids := make([]int, 0)
		counts += len(keys)
		for _, key := range keys {
			if !statics.AllNum(key[8:]) {
				counts--
				continue
			}
			id := 0
			fmt.Sscanf(key, "article:%d", &id)
			ids = append(ids, id)
		}
		DB.Model(&model.Article{}).Where("id in ?", ids).Find(&list)
		for _, i := range list {
			data, err := rdb.HGetAll(ctx, strconv.Itoa(int(i.ID))).Result()
			if err != nil {
				util.Log.Errorf("Failed to get hash data for key %d: %v", i.ID, err)
			}
			i.View = statics.StringToInt64(data["view"])
			i.Like = statics.StringToInt64(data["like"])
			i.Content = data["content"]
			i.Title = data["title"]
		}
		DB.Model(&model.Article{}).Save(&list)
		keyp.Put(keys)
		if cursor == 0 {
			break
		}
	}
	util.Log.Println("Sync Articles from Redis to SQLite finished, total:", counts)
	counts = 0
	cursor = 0
	for {
		keys := keyp.Get().([]string)
		var err error
		keys, cursor, err = rdb.Scan(ctx, cursor, "user:*", 10000).Result()
		if err != nil {
			util.Log.Errorf("Failed to scan keys from Redis: %v", err)
			return
		}
		var list []model.User
		ids := make([]int, 0)
		for _, key := range keys {
			id := 0
			fmt.Sscanf(key, "user:%d", &id)
			ids = append(ids, id)
		}
		DB.Model(&model.User{}).Where("id in ?", ids).Find(&list)
		for _, i := range list {
			data, err := rdb.HGetAll(ctx, strconv.Itoa(int(i.ID))).Result()
			if err != nil {
				util.Log.Errorf("Failed to get hash data for key %d: %v", i.ID, err)
			}
			i.Username = data["username"]
			i.Password = data["password"]
		}
		DB.Model(&model.User{}).Save(&list)
		counts += len(keys)
		keyp.Put(keys)
		if cursor == 0 {
			break
		}
	}
	util.Log.Println("Sync Users from Redis to SQLite finished, total:", counts)
}
func PullData(rdb *redis.Client, DB *gorm.DB) {
	ctx := context.Background()
	start := time.Now()
	util.Log.Println("Starting sync from SQLite to Redis...")
	var cursor uint64 = 0
	counts := 0
	for {
		var list []model.File
		DB.Model(&model.File{}).Offset(int(cursor)).Limit(200).Find(&list)
		for _, key := range list {
			rdb.HSet(ctx, fmt.Sprintf("file:%d", key.ID), "description", key.Description, "download_times", key.DownloadTimes, "file_path", key.FilePath, "author_id", key.AuthorID)
		}
		cursor += 200
		counts += len(list)
		if len(list) == 0 {
			break
		}
	}
	util.Log.Println("total files:", counts)
	cursor = 0
	counts = 0
	for {
		var list []model.Article
		DB.Model(&model.Article{}).Offset(int(cursor)).Limit(200).Find(&list)
		for _, key := range list {
			rdb.HSet(ctx, fmt.Sprintf("article:%d", key.ID), "view", key.View, "like", key.Like, "content", key.Content, "title", key.Title, "AuId", key.AuthorID, "cover", key.Cover)
		}
		cursor += 200
		counts += len(list)
		if len(list) == 0 {
			break
		}
	}
	util.Log.Println("total articles:", counts)
	cursor = 0
	counts = 0
	for {
		var list []model.User
		DB.Model(&model.User{}).Offset(int(cursor)).Limit(200).Find(&list)
		for _, key := range list {
			rdb.HSet(ctx, fmt.Sprintf("user:%d", key.ID), "username", key.Username, "password", key.Password, "priority", key.Priority)
		}
		cursor += 200
		counts += len(list)
		if len(list) == 0 {
			break
		}
	}
	util.Log.Println("total users:", counts)
	util.Log.Println("Sync Users from Sql to Redis finished ", time.Now().Sub(start).Seconds(), "seconds")
}
