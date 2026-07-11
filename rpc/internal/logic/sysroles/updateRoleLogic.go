package sysroleslogic

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/ent/sysrole"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"
	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
	"github.com/saas-zero/saas-zero-common/pkg/errno"
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
