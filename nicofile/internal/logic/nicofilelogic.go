package logic

import (
	"context"

	"main/nicofile/internal/svc"
	"main/nicofile/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type NicofileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewNicofileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NicofileLogic {
	return &NicofileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NicofileLogic) Nicofile(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
