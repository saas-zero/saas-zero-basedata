package syslogslogic

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLoginLogListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLoginLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLoginLogListLogic {
	return &GetLoginLogListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetLoginLogListLogic) GetLoginLogList(in *apps.LogPageReq) (*apps.LoginLogListResp, error) {
	// todo: add your logic here and delete this line

	return &apps.LoginLogListResp{}, nil
}
