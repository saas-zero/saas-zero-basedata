package sysapislogic

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetApiByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetApiByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetApiByIdLogic {
	return &GetApiByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetApiByIdLogic) GetApiById(in *apps.IdReq) (*apps.ApiResp, error) {
	// todo: add your logic here and delete this line

	return &apps.ApiResp{}, nil
}
