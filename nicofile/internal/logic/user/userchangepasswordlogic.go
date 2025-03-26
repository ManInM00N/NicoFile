package user

import (
	"context"
	"encoding/json"
	"fmt"
	"main/model"
	"main/pkg/encrypt"

	"main/nicofile/internal/svc"
	"main/nicofile/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserChangePasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserChangePasswordLogic {
	return &UserChangePasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserChangePasswordLogic) UserChangePassword(req *types.NewPasswordRequest) (resp *types.NewPasswordResponse, err error) {
	resp = &types.NewPasswordResponse{
		Error: false,
	}
	id, _ := l.ctx.Value("UserId").(json.Number).Int64()
	Password := encrypt.EncPassword(req.NewPassword)
	if err = l.svcCtx.DB.Model(model.User{}).Select("id,password").Where("id = ?", id).Update("Password", Password).Error; err != nil {
		resp.Error = true
		resp.Message = "密码修改失败"
	} else {
		resp.Message = "密码修改成功"
		err2 := l.svcCtx.Rdb.HSet(context.Background(), fmt.Sprintf("user:%d", id), "password", Password).Err()
		if err2 != nil {
			resp.Error = true
			resp.Message = "密码修改失败"
		}
	}
	return
}
