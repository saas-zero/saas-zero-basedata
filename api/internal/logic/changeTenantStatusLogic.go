// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/api/internal/svc"
	"github.com/saas-zero/saas-zero-basedata/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangeTenantStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangeTenantStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeTenantStatusLogic {
	return &ChangeTenantStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangeTenantStatusLogic) ChangeTenantStatus(req *types.TenantReq) (resp *types.BaseResp, err error) {
	// todo: add your logic here and delete this line

	return
}
