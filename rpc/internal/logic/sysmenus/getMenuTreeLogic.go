package sysmenuslogic

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/ent"
	"github.com/saas-zero/saas-zero-basedata/ent/sysmenu"
	"github.com/saas-zero/saas-zero-basedata/ent/sysrole"
	"github.com/saas-zero/saas-zero-basedata/ent/sysuser"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"
	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
	"github.com/saas-zero/saas-zero-common/pkg/errno"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetMenuTreeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMenuTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuTreeLogic {
	return &GetMenuTreeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMenuTreeLogic) GetMenuTree(_ *apps.EmptyReq) (*apps.MenuTreeResp, error) {
	userId := mixins.GetCurrentUserId(l.ctx)
	isAdmin := false
	var allMenus []*ent.SysMenu
	var err error

	if userId > 0 {
		user, uerr := l.svcCtx.DB.SysUser.ActiveQuery().
			Where(sysuser.IDEQ(userId)).
			WithRoles(func(q *ent.SysRoleQuery) {
				q.Where(sysrole.DeletedAtIsNil(), sysrole.StatusEQ(sysrole.StatusActive)).
					WithMenus(func(q *ent.SysMenuQuery) {
						q.Where(sysmenu.DeletedAtIsNil())
					})
			}).
			Only(l.ctx)
		if uerr == nil {
			for _, role := range user.Edges.Roles {
				if role.Code == "admin" {
					isAdmin = true
					break
				}
			}

			if !isAdmin {
				menuSet := make(map[int64]bool)
				for _, role := range user.Edges.Roles {
					for _, menu := range role.Edges.Menus {
						menuSet[menu.ID] = true
					}
				}
				if len(menuSet) > 0 {
					ids := make([]int64, 0, len(menuSet))
					for id := range menuSet {
						ids = append(ids, id)
					}
					allMenus, err = l.svcCtx.DB.SysMenu.ActiveQuery().
						Where(sysmenu.IDIn(ids...)).
						Order(ent.Asc(sysmenu.FieldSort)).
						All(l.ctx)
				}
			}
		}
	}

	if isAdmin {
		allMenus, err = l.svcCtx.DB.SysMenu.ActiveQuery().
			Order(ent.Asc(sysmenu.FieldSort)).
			All(l.ctx)
	}

	if err != nil {
		return nil, err
	}

	tree := buildMenuTree(allMenus, 0)
	return &apps.MenuTreeResp{
		Code: int32(errno.Success.Code),
		Msg:  errno.Success.Msg,
		Data: tree,
	}, nil
}
