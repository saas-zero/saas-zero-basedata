package sysroleslogic

import (
	"context"
	"fmt"

	"github.com/saas-zero/saas-zero-basedata/ent/sysrole"
	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRoleLogic {
	return &CreateRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateRoleLogic) CreateRole(in *apps.RoleReq) (*apps.RoleResp, error) {
	if in.Name == nil || *in.Name == "" {
		return &apps.RoleResp{Code: 500, Msg: "角色名称不能为空"}, nil
	}

	if in.Code == nil || *in.Code == "" {
		return &apps.RoleResp{Code: 500, Msg: "角色编码不能为空"}, nil
	}

	tenantId := mixins.GetCurrentTenantId(l.ctx)
	ctx := mixins.SetCurrentTenantId(l.ctx, tenantId)

	create := l.svcCtx.DB.SysRole.Create().
		SetName(*in.Name).
		SetCode(*in.Code)

	if in.Status != nil {
		create.SetStatus(sysrole.Status(*in.Status))
	}
	if in.Sort != nil {
		create.SetSort(uint32(*in.Sort))
	}
	if in.Remark != nil {
		create.SetRemark(*in.Remark)
	}
	if len(in.MenuIds) > 0 {
		create.AddMenuIDs(in.MenuIds...)
	}
	if len(in.ApiIds) > 0 {
		create.AddAPIIDs(in.ApiIds...)
	}

	role, err := create.Save(ctx)
	if err != nil {
		return &apps.RoleResp{Code: 500, Msg: fmt.Sprintf("创建角色失败: %v", err)}, nil
	}

	role, err = l.svcCtx.DB.SysRole.Query().
		Where(sysrole.IDEQ(role.ID)).
		WithMenus().
		WithApis().
		Only(l.ctx)
	if err != nil {
		return &apps.RoleResp{Code: 500, Msg: fmt.Sprintf("获取角色失败: %v", err)}, nil
	}

	return &apps.RoleResp{
		Code: 200,
		Msg:  "创建成功",
		Data: roleToPb(role, tenantId),
	}, nil
}