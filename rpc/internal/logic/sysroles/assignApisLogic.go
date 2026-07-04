package sysroleslogic

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AssignApisLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAssignApisLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignApisLogic {
	return &AssignApisLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AssignApisLogic) AssignApis(in *apps.RoleReq) (*apps.EmptyResp, error) {
	// todo: add your logic here and delete this line

	return &apps.EmptyResp{}, nil
}
