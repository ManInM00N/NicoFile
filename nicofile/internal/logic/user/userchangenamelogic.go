package user

import (
	"context"

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
	// todo: add your logic here and delete this line

	return
}
