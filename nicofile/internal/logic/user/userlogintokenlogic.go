package user

import (
	"context"
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
	if req.Id != l.ctx.Value("UserId").(int64) {
		resp.Error = true
		resp.Message = "身份认证失败"
	} else {
		var User model.User
		l.svcCtx.DB.Where("id = ?", req.Id).First(&User)
		resp.Username = User.Username
	}
	return
}
