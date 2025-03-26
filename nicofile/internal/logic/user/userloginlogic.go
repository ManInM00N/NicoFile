package user

import (
	"context"
	"fmt"
	"main/model"
	"main/pkg/encrypt"
	"main/pkg/jwt"

	"main/nicofile/internal/svc"
	"main/nicofile/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLoginLogic {
	return &UserLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLoginLogic) UserLogin(req *types.LoginRequest) (resp *types.AuthResponse, err error) {
	var User model.User
	l.svcCtx.DB.Model(&model.User{}).Select("id,username,password").Where("username = ?", req.Username).First(&User)
	if encrypt.EncPassword(req.Password) != User.Password {
		resp = &types.AuthResponse{
			Message: "账号或者密码错误",
		}
		return
	}
	l.svcCtx.Rdb.HSet(context.Background(), fmt.Sprintf("user:%d", User.ID), "username", User.Username, "priority", User.Priority, "password", User.Password)
	token, _ := jwt.BuildTokens(jwt.TokenOptions{AccessSecret: l.svcCtx.Config.Auth.AccessSecret, AccessExpire: l.svcCtx.Config.Auth.AccessExpire, Fields: map[string]interface{}{"UserId": User.ID}})
	resp = &types.AuthResponse{
		Message:  "登录成功",
		Token:    token.AccessToken,
		Username: User.Username,
	}
	return
}
