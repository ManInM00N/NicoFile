package user

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"google.golang.org/protobuf/proto"
	"main/model"
	"main/pkg/encrypt"
	"main/pkg/jwt"

	"main/nicofile/internal/svc"
	"main/nicofile/internal/types"
	"main/server/proto/kafka"

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
	resp = &types.AuthResponse{
		Error: false,
	}
	var User model.User
	User.Username = req.Username
	User.Password = encrypt.EncPassword(req.Password)
	err2 := l.svcCtx.DB.Create(&User).Error
	if err2 != nil {
		resp.Error = true
		resp.Message = "名称非法或已存在"
		return
	}
	event := &kafka.UserMonitor{
		Message: "A new user has been registered",
		Warning: false,
		UserId:  uint32(User.ID),
	}
	data, _ := proto.Marshal(event)
	(*l.svcCtx.Producer).Input() <- &sarama.ProducerMessage{
		Topic: "data-monitor-test", // l.svcCtx.Config.Kafka.Topic,
		Value: sarama.ByteEncoder(data),
	}
	l.svcCtx.Rdb.HSet(context.Background(), fmt.Sprintf("user:%d", User.ID), "username", User.Username, "priority", User.Priority, "password", User.Password)
	token, _ := jwt.BuildTokens(jwt.TokenOptions{AccessSecret: l.svcCtx.Config.Auth.AccessSecret, AccessExpire: l.svcCtx.Config.Auth.AccessExpire, Fields: map[string]interface{}{"UserId": User.ID}})
	resp.Message = "注册成功"
	resp.Token = token.AccessToken
	resp.Username = User.Username
	return
}
