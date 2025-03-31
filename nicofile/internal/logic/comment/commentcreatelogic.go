package comment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"main/model"
	"main/nicofile/internal/svc"
	"main/nicofile/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommentCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentCreateLogic {
	return &CommentCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentCreateLogic) CommentCreate(req *types.CommentCreateRequest) (resp *types.CommentResponse, err error) {
	resp = &types.CommentResponse{
		Error: false,
	}
	id, _ := l.ctx.Value("UserId").(json.Number).Int64()
	//ip, _, err := net.SplitHostPort(l.ctx.Value("remoteAddr").(string))
	//if err != nil {
	ip := "unknown"
	//}
	data := model.Comment{
		Content:   req.Content,
		IP:        ip,
		ArticleID: uint(req.ArticleId),
		AuthorId:  uint(id),
		Status:    "pending",
		ParentID:  nil,
	}
	if req.ParentId != 0 {
		data.ParentID = new(uint)
		*data.ParentID = uint(req.ParentId)
	}
	if err2 := l.svcCtx.DB.Model(&model.Comment{}).Create(&data).Error; err2 != nil {
		fmt.Println(err2, id, req.ArticleId)
		return nil, errors.New("创建失败error creating")
	}
	return
}
