package sysroleslogic

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/ent/sysrole"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"
	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetRoleByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRoleByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoleByIdLogic {
	return &GetRoleByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetRoleByIdLogic) GetRoleById(in *apps.IdReq) (*apps.RoleResp, error) {
	tenantId := mixins.GetCurrentTenantId(l.ctx)

	r, err := l.svcCtx.DB.SysRole.TenantQuery(tenantId).
		Where(sysrole.IDEQ(in.GetId())).
		WithMenus().
		Only(l.ctx)
	if err != nil {
		return nil, err
	}
	resp := roleToResp(r)
	resp.ApiIds = roleApiIds(l.svcCtx.Enforcer, r.Code, tenantId)
	return &apps.RoleResp{
		Code: 200,
		Msg:  "success",
		Data: resp,
	}, nil
}
