package handler

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"main/config"
	"main/model"
	"main/pkg/jwt"
	"main/server/proto/auth"
)

type AuthServiceServer struct {
	auth.UnimplementedAuthServiceServer
}

func (s *AuthServiceServer) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	var user model.User
	result := config.DB.Where("username = ?", req.Username).First(&user)

	if result.Error != nil {
		return &auth.LoginResponse{
			Success: false,
			Message: "User not found",
		}, nil
	}
	// 比对密码
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return &auth.LoginResponse{
			Success: false,
			Message: "Invalid password",
		}, nil
	}

	return &auth.LoginResponse{
		Success: true,
		Message: "Login successful",
	}, nil
}

func (s *AuthServiceServer) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	var user model.User
	result := config.DB.Where("username = ?", req.Username).First(&user)
	user.Username = req.Username
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return &auth.RegisterResponse{
			Success: false,
			Message: "User name existed",
			Cookie:  "",
		}, nil
	}
	res, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user.Password = string(res)
	config.DB.Create(&user)
	fmt.Println(user.ID)
	cookie, _ := jwt.BuildTokens(
		jwt.TokenOptions{
			AccessSecret: "114514",
			AccessExpire: 3600,
			Fields:       map[string]interface{}{"UserId": user.ID},
		})
	return &auth.RegisterResponse{
		Success: true,
		Message: "Register Successful",
		Cookie:  cookie.AccessToken,
	}, nil
}
