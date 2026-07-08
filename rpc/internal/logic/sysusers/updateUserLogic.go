package sysuserslogic

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/ent/sysuser"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"
	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
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

	result, err := update.Save(ctx)
	if err != nil {
		return nil, err
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

	return &apps.UserResp{Code: 200, Msg: "success", Data: userToResp(u)}, nil
}
