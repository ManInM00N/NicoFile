package user

import (
	"context"
	"encoding/json"
	"main/model"

	"main/nicofile/internal/svc"
	"main/nicofile/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserLoginTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLoginTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLoginTokenLogic {
	return &UserLoginTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLoginTokenLogic) UserLoginToken(req *types.LoginTokenRequest) (resp *types.CheckResponse, err error) {
	resp = &types.CheckResponse{
		Error:    false,
		Message:  "身份认证成功",
		Username: "",
	}
	id, _ := l.ctx.Value("UserId").(json.Number).Int64()
	var User model.User
	if err = l.svcCtx.DB.Model(&model.User{}).Select("id,username").Where("id = ?", id).First(&User).Error; err != nil {
		resp.Error = true
		resp.Message = "身份认证失败"
		return
	}
	resp.Username = User.Username
	return
}
