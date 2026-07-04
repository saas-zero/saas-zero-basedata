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

type GetMenuRoutersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMenuRoutersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuRoutersLogic {
	return &GetMenuRoutersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMenuRoutersLogic) GetMenuRouters(in *apps.EmptyReq) (*apps.MenuTreeResp, error) {
	tenantId := mixins.GetCurrentTenantId(l.ctx)

	menus, err := l.svcCtx.DB.SysMenu.TenantQuery(tenantId).
		Where(sysmenu.StatusEQ(sysmenu.StatusActive)).
		Order(ent.Asc(sysmenu.FieldSort)).
		All(l.ctx)
	if err != nil {
		return &apps.MenuTreeResp{Code: 500, Msg: fmt.Sprintf("获取路由菜单失败: %v", err)}, nil
	}

	tree := buildMenuTree(menus, 0, tenantId)

	return &apps.MenuTreeResp{
		Code: 200,
		Msg:  "success",
		Data: tree,
	}, nil
}