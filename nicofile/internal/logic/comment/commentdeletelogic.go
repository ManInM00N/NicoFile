package comment

import (
	"context"
	"errors"
	"fmt"
	"main/model"

	"main/nicofile/internal/svc"
	"main/nicofile/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommentDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentDeleteLogic {
	return &CommentDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentDeleteLogic) CommentDelete(req *types.DownloadIMGRequest) (resp *types.Response, err error) {
	resp = &types.Response{
		Message: "NULL",
	}
	if err2 := l.svcCtx.DB.Unscoped().Model(&model.Comment{}).Where(" parent_id = ?", req.ID).Delete(&model.Comment{}).Error; err2 != nil {
		fmt.Println(err2, req.ID)

		return nil, errors.New("删除失败error deleting")
	} else {
		if err2 := l.svcCtx.DB.Unscoped().Model(&model.Comment{}).Where("id = ?", req.ID).Delete(&model.Comment{}).Error; err2 != nil {
			fmt.Println(err2, req.ID)
			return nil, errors.New("删除失败error deleting")
		}
	}
	return
}
