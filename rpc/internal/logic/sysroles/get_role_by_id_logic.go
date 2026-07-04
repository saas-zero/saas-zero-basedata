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

	role, err := l.svcCtx.DB.SysRole.Query().
		Where(sysrole.IDEQ(int64(in.Id))).
		WithMenus().
		WithApis().
		Only(l.ctx)
	if err != nil {
		return &apps.RoleResp{Code: 500, Msg: fmt.Sprintf("获取角色失败: %v", err)}, nil
	}

	return &apps.RoleResp{
		Code: 200,
		Msg:  "success",
		Data: roleToPb(role, tenantId),
	}, nil
}

func roleToPb(r *ent.SysRole, tenantId int64) *apps.Role {
	sort := int32(r.Sort)
	isSystem := r.IsSystem
	role := &apps.Role{
		Id:       &r.ID,
		IdStr:    strPtr(fmt.Sprintf("%d", r.ID)),
		Name:     &r.Name,
		Code:     &r.Code,
		IsSystem: &isSystem,
		Status:   strPtr(string(r.Status)),
		Sort:     &sort,
		Remark:   &r.Remark,
		TenantId: &tenantId,
		TenantIdStr: strPtr(fmt.Sprintf("%d", tenantId)),
	}

	if len(r.Edges.Menus) > 0 {
		menuIds := make([]int64, 0, len(r.Edges.Menus))
		for _, m := range r.Edges.Menus {
			menuIds = append(menuIds, m.ID)
		}
		role.MenuIds = menuIds
	}

	if len(r.Edges.Apis) > 0 {
		apiIds := make([]int64, 0, len(r.Edges.Apis))
		for _, a := range r.Edges.Apis {
			apiIds = append(apiIds, a.ID)
		}
		role.ApiIds = apiIds
	}

	createdAt := r.CreatedAt.Unix()
	role.CreatedAt = &createdAt
	role.CreatedBy = &r.CreatedBy

	updatedAt := r.UpdatedAt.Unix()
	role.UpdatedAt = &updatedAt
	role.UpdatedBy = &r.UpdatedBy

	return role
}

func strPtr(s string) *string {
	return &s
}