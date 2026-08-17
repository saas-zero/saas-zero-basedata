package sysroleslogic

import (
	"context"
	"fmt"

	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"
	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
	"github.com/saas-zero/saas-zero-common/pkg/errno"
	"github.com/zeromicro/go-zero/core/logx"
)

type AssignMenusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAssignMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignMenusLogic {
	return &AssignMenusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AssignMenusLogic) AssignMenus(in *apps.RoleReq) (*apps.EmptyResp, error) {
	userId := mixins.GetCurrentUserId(l.ctx)
	userName := mixins.GetCurrentUserName(l.ctx)
	ctx := mixins.SetCurrentUserId(l.ctx, userId)
	ctx = mixins.SetCurrentUserName(ctx, userName)

	// 系统内置角色不可改菜单权限
	role, err := l.svcCtx.DB.SysRole.Get(ctx, in.GetId())
	if err != nil {
		return nil, err
	}
	if role.IsSystem {
		return nil, errno.New(errno.InvalidParam.Code, fmt.Sprintf("系统内置角色「%s」不可修改", role.Name))
	}

	// 继承式授权：只能把当前用户自己拥有的菜单授给别人
	if err := checkAssignableMenus(l.svcCtx, l.ctx, in.GetMenuIds()); err != nil {
		return nil, err
	}

	err = l.svcCtx.DB.SysRole.UpdateOneID(in.GetId()).
		ClearMenus().
		AddMenuIDs(in.GetMenuIds()...).
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	return &apps.EmptyResp{Code: int32(errno.Success.Code), Msg: errno.Success.Msg}, nil
}
