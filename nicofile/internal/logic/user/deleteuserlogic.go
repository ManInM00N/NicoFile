package user

import (
	"context"
	"encoding/json"
	"main/model"

	"main/nicofile/internal/svc"
	"main/nicofile/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteUserLogic) DeleteUser(req *types.DeleteUserRequest) (resp *types.DeleteUserResponse, err error) {
	id, _ := l.ctx.Value("UserId").(json.Number).Int64()
	resp = &types.DeleteUserResponse{
		Error: false,
	}
	if id != req.Userid && id != 0 {
		resp.Message = "无权删除"
		resp.Error = true
		return
	}
	if err2 := l.svcCtx.DB.Unscoped().Delete(&model.User{}, req.Userid).Error; err2 != nil {
		resp.Message = "删除失败或者用户不存在"
		resp.Error = true
	}

	return
}
