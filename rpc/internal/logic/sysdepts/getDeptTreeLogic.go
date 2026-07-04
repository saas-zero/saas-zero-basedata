package sysdeptslogic

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDeptTreeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDeptTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDeptTreeLogic {
	return &GetDeptTreeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetDeptTreeLogic) GetDeptTree(in *apps.EmptyReq) (*apps.DeptTreeResp, error) {
	// todo: add your logic here and delete this line

	return &apps.DeptTreeResp{}, nil
}
