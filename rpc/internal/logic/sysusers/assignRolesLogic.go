package sysuserslogic

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"
	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
	"github.com/saas-zero/saas-zero-common/pkg/errno"
)

type AssignRolesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAssignRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignRolesLogic {
	return &AssignRolesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AssignRolesLogic) AssignRoles(in *apps.UserReq) (*apps.EmptyResp, error) {
	userId := mixins.GetCurrentUserId(l.ctx)
	userName := mixins.GetCurrentUserName(l.ctx)
	ctx := mixins.SetCurrentUserId(l.ctx, userId)
	ctx = mixins.SetCurrentUserName(ctx, userName)

	err := l.svcCtx.DB.SysUser.UpdateOneID(in.GetId()).
		ClearRoles().
		AddRoleIDs(in.GetRoleIds()...).
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	l.svcCtx.Redis.Incr(fmt.Sprintf("token_version:%d", in.GetId()))
	return &apps.EmptyResp{Code: int32(errno.Success.Code), Msg: errno.Success.Msg}, nil
}
