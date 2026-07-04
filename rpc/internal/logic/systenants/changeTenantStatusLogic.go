package systenantslogic

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangeTenantStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChangeTenantStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeTenantStatusLogic {
	return &ChangeTenantStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ChangeTenantStatusLogic) ChangeTenantStatus(in *apps.TenantReq) (*apps.EmptyResp, error) {
	// todo: add your logic here and delete this line

	return &apps.EmptyResp{}, nil
}
