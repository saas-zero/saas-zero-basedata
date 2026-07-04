package sysmenuslogic

import (
	"context"
	"fmt"

	"github.com/saas-zero/saas-zero-basedata/ent"
	"github.com/saas-zero/saas-zero-basedata/ent/sysmenu"
	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

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

func (l *GetMenuTreeLogic) GetMenuTree(in *apps.EmptyReq) (*apps.MenuTreeResp, error) {
	tenantId := mixins.GetCurrentTenantId(l.ctx)

	menus, err := l.svcCtx.DB.SysMenu.TenantQuery(tenantId).
		Order(ent.Asc(sysmenu.FieldSort)).
		All(l.ctx)
	if err != nil {
		return &apps.MenuTreeResp{Code: 500, Msg: fmt.Sprintf("查询菜单树失败: %v", err)}, nil
	}

	tree := buildMenuTree(menus, 0, tenantId)

	return &apps.MenuTreeResp{
		Code: 200,
		Msg:  "success",
		Data: tree,
	}, nil
}

func buildMenuTree(menus []*ent.SysMenu, parentId int64, tenantId int64) []*apps.Menu {
	tree := make([]*apps.Menu, 0)
	for _, m := range menus {
		if m.ParentID == parentId {
			menu := menuToPb(m, tenantId)
			menu.Children = buildMenuTree(menus, m.ID, tenantId)
			tree = append(tree, menu)
		}
	}
	return tree
}