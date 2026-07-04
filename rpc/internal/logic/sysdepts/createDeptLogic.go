package sysdeptslogic

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateDeptLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateDeptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDeptLogic {
	return &CreateDeptLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateDeptLogic) CreateDept(in *apps.DeptReq) (*apps.DeptResp, error) {
	// todo: add your logic here and delete this line

	return &apps.DeptResp{}, nil
}
