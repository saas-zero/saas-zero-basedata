package sysroleslogic

import (
	"context"
	"fmt"

	"github.com/saas-zero/saas-zero-basedata/ent"
	"github.com/saas-zero/saas-zero-basedata/ent/sysrole"
	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

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

	query := l.svcCtx.DB.SysRole.TenantQuery(tenantId).
		WithMenus().
		WithApis()

	if in.Name != nil && *in.Name != "" {
		query.Where(sysrole.NameContains(*in.Name))
	}
	if in.Code != nil && *in.Code != "" {
		query.Where(sysrole.CodeContains(*in.Code))
	}
	if in.Status != nil && *in.Status != "" {
		query.Where(sysrole.StatusEQ(sysrole.Status(*in.Status)))
	}

	total, err := query.Count(l.ctx)
	if err != nil {
		return &apps.RoleListResp{Code: 500, Msg: fmt.Sprintf("查询角色总数失败: %v", err)}, nil
	}

	page := int(in.Page)
	size := int(in.Size)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	roles, err := query.
		Offset((page - 1) * size).
		Limit(size).
		Order(ent.Asc(sysrole.FieldSort)).
		All(l.ctx)
	if err != nil {
		return &apps.RoleListResp{Code: 500, Msg: fmt.Sprintf("查询角色列表失败: %v", err)}, nil
	}

	list := make([]*apps.Role, 0, len(roles))
	for _, r := range roles {
		list = append(list, roleToPb(r, tenantId))
	}

	return &apps.RoleListResp{
		Code:  200,
		Msg:   "success",
		List:  list,
		Total: int64(total),
	}, nil
}