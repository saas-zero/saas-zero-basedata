package sysroleslogic

import (
	"context"

	"fmt"
	"github.com/saas-zero/saas-zero-basedata/ent/sysrole"
	"github.com/saas-zero/saas-zero-basedata/ent/sysuser"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"
	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
	"github.com/saas-zero/saas-zero-common/pkg/errno"
	"github.com/saas-zero/saas-zero-common/pkg/id"
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
	userId := mixins.GetCurrentUserId(l.ctx)
	userName := mixins.GetCurrentUserName(l.ctx)
	ctx := mixins.SetCurrentUserId(l.ctx, userId)
	ctx = mixins.SetCurrentUserName(ctx, userName)

	// 记录变更前信息：旧 code（用于同步 Casbin）、旧状态（用于禁用踢出）
	oldRole, err := l.svcCtx.DB.SysRole.Get(ctx, in.GetId())
	if err != nil {
		return nil, err
	}
	oldCode := oldRole.Code
	oldStatus := oldRole.Status

	// 系统内置角色（is_system=true）不可修改
	if oldRole.IsSystem {
		return nil, errno.New(errno.InvalidParam.Code, fmt.Sprintf("系统内置角色「%s」不可修改", oldRole.Name))
	}

	update := l.svcCtx.DB.SysRole.UpdateOneID(in.GetId())
	if in.Name != nil {
		update.SetName(in.GetName())
	}
	if in.Code != nil {
		update.SetCode(in.GetCode())
	}
	if in.Status != nil {
		update.SetStatus(sysrole.Status(in.GetStatus()))
	}
	if in.Sort != nil {
		update.SetSort(uint32(in.GetSort()))
	}
	if in.Remark != nil {
		update.SetRemark(in.GetRemark())
	}

	result, err := update.Save(ctx)
	if err != nil {
		return nil, err
	}

	// 修改角色 code：同步 Casbin 策略 sub（删旧 code 策略，重建为新 code）
	if in.Code != nil && in.GetCode() != oldCode {
		dom := id.ToString(mixins.GetCurrentTenantId(ctx))
		if l.svcCtx.Enforcer != nil {
			oldPolicies, _ := l.svcCtx.Enforcer.GetFilteredPolicy(0, oldCode, dom)
			l.svcCtx.Enforcer.RemoveFilteredPolicy(0, oldCode, dom)
			for _, p := range oldPolicies {
				if len(p) >= 5 {
					l.svcCtx.Enforcer.AddPolicy(in.GetCode(), dom, p[2], p[3], p[4])
				}
			}
		}
	}

	// 禁用角色（active → inactive）：踢掉拥有该角色的所有用户
	if in.Status != nil && sysrole.Status(in.GetStatus()) == sysrole.StatusInactive && oldStatus == sysrole.StatusActive {
		users, err := l.svcCtx.DB.SysUser.Query().
			Where(sysuser.HasRolesWith(sysrole.IDEQ(result.ID))).
			All(ctx)
		if err == nil {
			for _, u := range users {
				l.svcCtx.Redis.Incr(fmt.Sprintf("token_version:%d", u.ID))
			}
		}
	}

	if len(in.GetMenuIds()) > 0 {
		l.svcCtx.DB.SysRole.UpdateOneID(result.ID).ClearMenus().AddMenuIDs(in.GetMenuIds()...).Exec(ctx)
	}
	r, err := l.svcCtx.DB.SysRole.Query().Where(sysrole.IDEQ(result.ID)).WithMenus().Only(ctx)
	if err != nil {
		return nil, err
	}
	return &apps.RoleResp{
		Code: int32(errno.Success.Code),
		Msg:  errno.Success.Msg,
		Data: roleToResp(r),
	}, nil
}
