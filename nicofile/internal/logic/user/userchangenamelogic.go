package user

import (
	"context"
	"encoding/json"
	"fmt"
	"main/model"

	"main/nicofile/internal/svc"
	"main/nicofile/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserChangeNameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserChangeNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserChangeNameLogic {
	return &UserChangeNameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserChangeNameLogic) UserChangeName(req *types.NewNameRequest) (resp *types.NewNameResponse, err error) {
	resp = &types.NewNameResponse{
		Error: false,
	}
	id, _ := l.ctx.Value("UserId").(json.Number).Int64()
	var User model.User
	if l.svcCtx.DB.Model(model.User{}).Select("id,username").Where("username = ?", req.NewName).First(&User).Error == nil && User.ID != uint(id) {
		resp.Error = true
		resp.Message = "用户名已存在"
		return
	}
	if err = l.svcCtx.DB.Model(model.User{}).Select("id,username").Where("id = ?", id).Update("Username", req.NewName).Error; err != nil {
		resp.Error = true
		resp.Message = "用户名修改失败"
		return
	} else {
		l.svcCtx.Rdb.HSet(context.Background(), fmt.Sprintf("user:%d", id), "username", req.NewName).Err()
		//l.svcCtx.Rdb.Expire(context.Background(), fmt.Sprintf("user:%d", id), 60*60).Err()
	}
	return
}
