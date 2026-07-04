package sysroleslogic

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/ent"
	"github.com/saas-zero/saas-zero-basedata/ent/sysrole"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"
	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetRoleListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRoleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoleListLogic {
	return &GetRoleListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetRoleListLogic) GetRoleList(in *apps.RolePageReq) (*apps.RoleListResp, error) {
	tenantId := mixins.GetCurrentTenantId(l.ctx)

	query := l.svcCtx.DB.SysRole.TenantQuery(tenantId)
	if in.GetName() != "" {
		query = query.Where(sysrole.NameContains(in.GetName()))
	}
	if in.GetCode() != "" {
		query = query.Where(sysrole.CodeContains(in.GetCode()))
	}
	if in.GetStatus() != "" {
		query = query.Where(sysrole.StatusEQ(sysrole.Status(in.GetStatus())))
	}

	total, err := query.Count(l.ctx)
	if err != nil {
		return nil, err
	}

	page := int(in.GetPage())
	size := int(in.GetSize())
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	roles, err := query.
		Offset((page - 1) * size).
		Limit(size).
		Order(ent.Asc(sysrole.FieldSort)).
		WithMenus().
		All(l.ctx)
	if err != nil {
		return nil, err
	}

	list := make([]*apps.Role, len(roles))
	for i, r := range roles {
		list[i] = roleToResp(r)
	}
	return &apps.RoleListResp{
		Code:  200,
		Msg:   "success",
		List:  list,
		Total: int64(total),
	}, nil
}
