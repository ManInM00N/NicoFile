package file

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"main/model"
	"main/nicofile/internal/svc"
	"main/nicofile/internal/types"
	"net/http"
)

type FileDownloadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileDownloadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileDownloadLogic {
	return &FileDownloadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileDownloadLogic) FileDownload(req *types.FileDownloadRequest, w http.ResponseWriter, file model.File) (resp *types.FileDownloadResponse, err error) {
	if !l.svcCtx.Config.Redis.Disabled {
		key := fmt.Sprintf("file:%d", file.ID)
		l.svcCtx.Rdb.HIncrBy(context.Background(), key, "download_times", 1)
	} else {
		l.svcCtx.DB.Model(&model.File{}).Where("file_path = ?", req.Url).UpdateColumn("download_times", gorm.Expr("download_times + ?", 1))
	}
	//resp = &types.FileDownloadResponse{}
	//f, err := os.OpenFile(l.svcCtx.Config.StoragePath+"/"+file.FilePath, os.O_RDONLY, 0666)
	////stat, _ := f.Stat()
	////w.Header().Set("Content-Length", strconv.FormatInt(stat.Size(), 10))
	//if err != nil {
	//	return nil, err
	//}
	//defer f.Close()
	//r := bufio.NewReader(f)
	//buffer := util.Pool.Get().([]byte)
	////limiter := ratelimit.New(1024)
	//for {
	//	//limiter.Take()
	//	_, err = r.Read(buffer)
	//	if err != nil {
	//		if err == io.EOF {
	//			break
	//		}
	//		return nil, err
	//	}
	//	w.Write(buffer)
	//
	//}
	return
}
