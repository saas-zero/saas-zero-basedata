package sysmenuslogic

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMenuRoutersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMenuRoutersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuRoutersLogic {
	return &GetMenuRoutersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMenuRoutersLogic) GetMenuRouters(in *apps.EmptyReq) (*apps.MenuTreeResp, error) {
	// todo: add your logic here and delete this line

	return &apps.MenuTreeResp{}, nil
}
