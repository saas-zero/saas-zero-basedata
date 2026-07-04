package sysdeptslogic

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDeptListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDeptListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDeptListLogic {
	return &GetDeptListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetDeptListLogic) GetDeptList(in *apps.DeptPageReq) (*apps.DeptListResp, error) {
	// todo: add your logic here and delete this line

	return &apps.DeptListResp{}, nil
}
