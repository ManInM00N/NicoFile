package user

import (
	"context"
	"main/model"
	"main/pkg/encrypt"
	"main/pkg/jwt"

	"main/nicofile/internal/svc"
	"main/nicofile/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRegisterLogic {
	return &UserRegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserRegisterLogic) UserRegister(req *types.RegisterRequest) (resp *types.AuthResponse, err error) {
	// todo: add your logic here and delete this line
	var User model.User
	User.Username = req.Username
	User.Password = encrypt.EncPassword(req.Password)
	l.svcCtx.DB.Create(&User)
	token, _ := jwt.BuildTokens(jwt.TokenOptions{AccessSecret: l.svcCtx.Config.Auth.AccessSecret, AccessExpire: l.svcCtx.Config.Auth.AccessExpire, Fields: map[string]interface{}{"UserId": User.ID}})
	resp = &types.AuthResponse{
		Message: "注册成功",
		Token:   token.AccessToken,
	}
	return
}
