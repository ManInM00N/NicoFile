
<p align="center">

<img style="width :256px" src="https://raw.githubusercontent.com/Dip-a-scent-of-Blossom/NicoFile_Web/refs/heads/main/public/logo.png">
<br>
<a href="https://golang.google.cn/dl/"> <img src="https://img.shields.io/github/go-mod/go-version/ManInM00N/Go_Pixiv"> </a>
<a href="https://github.com/ManInM00N/Go_Pixiv/blob/master/LICENSE"><img src="https://img.shields.io/github/license/ManInM00N/NicoFile"> </a>
<a href="https://github.com/zeromicro/go-zero"> <img src="https://img.shields.io/badge/go--zero-v1.8.1-red"> </a>

</p>

### 此项目致力于打造一个开箱即用的基于go-zero的文件管理服务,集成了博客系统，支持文件上传、下载、删除、预览、文章点赞收藏等功能。

- [x] 总体文件功能 :tada:
  - [x] 断点续传/文件分片功能 :tada:
  - [x] 文件上传功能 :tada:
  - [ ] git同步功能 :tada:
  - [x] 优化预览缓存机制 :wave:
- [x] 文章基础功能 :monocle_face:
  - [x] 热点文章排行榜 :monocle_face:
  - [x] 图库 
- [x] 引入ES
  - [x] 文章搜索功能 :monocle_face:
  - [x] 数据自动迁移功能 :monocle_face:
- [ ] 用户群组优化
- [ ] Redis 集群化
- [x] kafka集成prometheus监控 :alien:
  - [x] grafana分析流量 :alien:
- [x] 引入etcd管理微服务 :zap:
  - [x] 迁移用户微服务
  - [x] 迁移排行榜 
- [ ] 引入gpt-3.5 接口 

前端部分请移步：[NicoFile_Web](https://github.com/Dip-a-scent-of-Blossom/NicoFile_Web)


To QuickStart ：
```shell
docker-compose up -d
go run nicofile/nicofile.go
go run server/server.go 
```
To generate api:
```shell
goctl api go --api=./nicofile/nicofile.api --dir=./nicofile 
protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. server/proto/*/*.proto
```