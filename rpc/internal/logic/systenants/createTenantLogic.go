package systenantslogic

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateTenantLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateTenantLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTenantLogic {
	return &CreateTenantLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateTenantLogic) CreateTenant(in *apps.TenantReq) (*apps.TenantResp, error) {
	// todo: add your logic here and delete this line

	return &apps.TenantResp{}, nil
}
