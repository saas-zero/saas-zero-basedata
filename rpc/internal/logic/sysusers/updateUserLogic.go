package sysuserslogic

import (
	"context"
	"fmt"

	"github.com/saas-zero/saas-zero-basedata/ent/sysuser"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"
	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
	"github.com/saas-zero/saas-zero-common/pkg/errno"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserLogic) UpdateUser(in *apps.UserReq) (*apps.UserResp, error) {
	userId := mixins.GetCurrentUserId(l.ctx)
	userName := mixins.GetCurrentUserName(l.ctx)
	ctx := mixins.SetCurrentUserId(l.ctx, userId)
	ctx = mixins.SetCurrentUserName(ctx, userName)

	update := l.svcCtx.DB.SysUser.UpdateOneID(in.GetId())
	if in.Nickname != nil {
		update.SetNickname(in.GetNickname())
	}
	if in.Mobile != nil {
		update.SetMobile(in.GetMobile())
	}
	if in.Email != nil {
		update.SetEmail(in.GetEmail())
	}
	if in.DeptId != nil {
		update.SetDeptID(in.GetDeptId())
	}
	if in.Status != nil {
		update.SetStatus(sysuser.Status(in.GetStatus()))
	}
	if in.Remark != nil {
		update.SetRemark(in.GetRemark())
	}

	// 记录更新前的状态，用于判断是否发生"禁用"
	var prevStatus sysuser.Status
	if in.Status != nil {
		prev, err := l.svcCtx.DB.SysUser.Get(ctx, in.GetId())
		if err != nil {
			return nil, err
		}
		prevStatus = prev.Status
	}

	result, err := update.Save(ctx)
	if err != nil {
		return nil, err
	}

	// 禁用用户（active → inactive）时递增 token_version，使其所有旧 token 立即失效，
	// 被禁用的用户无法继续访问系统。
	if in.Status != nil &&
		prevStatus == sysuser.StatusActive &&
		sysuser.Status(in.GetStatus()) == sysuser.StatusInactive {
		l.svcCtx.Redis.Incr(fmt.Sprintf("token_version:%d", in.GetId()))
	}

	if len(in.GetRoleIds()) > 0 {
		l.svcCtx.DB.SysUser.UpdateOneID(result.ID).
			ClearRoles().
			AddRoleIDs(in.GetRoleIds()...).
			Exec(ctx)
	}

	u, err := l.svcCtx.DB.SysUser.Query().
		Where(sysuser.IDEQ(result.ID)).
		WithRoles().
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return &apps.UserResp{Code: int32(errno.Success.Code), Msg: errno.Success.Msg, Data: userToResp(u)}, nil
}
