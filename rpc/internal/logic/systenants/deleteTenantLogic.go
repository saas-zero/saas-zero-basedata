package systenantslogic

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteTenantLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteTenantLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteTenantLogic {
	return &DeleteTenantLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteTenantLogic) DeleteTenant(in *apps.IdsReq) (*apps.EmptyResp, error) {
	// todo: add your logic here and delete this line

	return &apps.EmptyResp{}, nil
}
