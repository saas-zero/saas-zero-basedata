package sysapislogic

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateApiLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateApiLogic {
	return &UpdateApiLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateApiLogic) UpdateApi(in *apps.ApiReq) (*apps.ApiResp, error) {
	// todo: add your logic here and delete this line

	return &apps.ApiResp{}, nil
}
