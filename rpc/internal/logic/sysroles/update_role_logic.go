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

type UpdateRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRoleLogic {
	return &UpdateRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateRoleLogic) UpdateRole(in *apps.RoleReq) (*apps.RoleResp, error) {
	if in.Id == nil || *in.Id <= 0 {
		return &apps.RoleResp{Code: 500, Msg: "角色ID不能为空"}, nil
	}

	tenantId := mixins.GetCurrentTenantId(l.ctx)

	update := l.svcCtx.DB.SysRole.UpdateOneID(int64(*in.Id))

	if in.Name != nil {
		update.SetName(*in.Name)
	}
	if in.Code != nil {
		update.SetCode(*in.Code)
	}
	if in.Status != nil {
		update.SetStatus(sysrole.Status(*in.Status))
	}
	if in.Sort != nil {
		update.SetSort(uint32(*in.Sort))
	}
	if in.Remark != nil {
		update.SetRemark(*in.Remark)
	}

	role, err := update.Save(l.ctx)
	if err != nil {
		return &apps.RoleResp{Code: 500, Msg: fmt.Sprintf("更新角色失败: %v", err)}, nil
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
		Msg:  "更新成功",
		Data: roleToPb(role, tenantId),
	}, nil
}