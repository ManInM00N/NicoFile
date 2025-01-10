package handler

import (
	"context"
	"fmt"
	"main/config"
	"main/model"
	"main/server/proto"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthServiceServer struct {
	proto.UnimplementedAuthServiceServer
}

func (s *AuthServiceServer) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	var user model.User
	result := config.DB.Where("username = ?", req.Username).First(&user)

	if result.Error != nil {
		return &proto.LoginResponse{
			Success: false,
			Message: "User not found",
		}, nil
	}
	// 比对密码
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return &proto.LoginResponse{
			Success: false,
			Message: "Invalid password",
		}, nil
	}

	return &proto.LoginResponse{
		Success: true,
		Message: "Login successful",
	}, nil
}

func (s *AuthServiceServer) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	var user model.User
	result := config.DB.Where("username = ?", req.Username).First(&user)
	user.Username = req.Username
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return &proto.RegisterResponse{
			Success: false,
			Message: "User name existed",
			Cookie:  "",
		}, nil
	}
	res, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user.Password = string(res)
	config.DB.Create(&user)
	fmt.Println(user.ID)
	cookie, _ := GenerateToken(user.ID)
	return &proto.RegisterResponse{
		Success: true,
		Message: "Register Successful",
		Cookie:  cookie,
	}, nil
}
