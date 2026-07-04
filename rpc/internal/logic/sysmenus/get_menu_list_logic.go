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

type GetMenuListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMenuListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuListLogic {
	return &GetMenuListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMenuListLogic) GetMenuList(in *apps.MenuPageReq) (*apps.MenuListResp, error) {
	tenantId := mixins.GetCurrentTenantId(l.ctx)

	query := l.svcCtx.DB.SysMenu.TenantQuery(tenantId)

	if in.Name != nil && *in.Name != "" {
		query.Where(sysmenu.NameContains(*in.Name))
	}
	if in.Status != nil && *in.Status != "" {
		query.Where(sysmenu.StatusEQ(sysmenu.Status(*in.Status)))
	}

	total, err := query.Count(l.ctx)
	if err != nil {
		return &apps.MenuListResp{Code: 500, Msg: fmt.Sprintf("查询菜单总数失败: %v", err)}, nil
	}

	page := int(in.Page)
	size := int(in.Size)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	menus, err := query.
		Offset((page - 1) * size).
		Limit(size).
		Order(ent.Asc(sysmenu.FieldSort)).
		All(l.ctx)
	if err != nil {
		return &apps.MenuListResp{Code: 500, Msg: fmt.Sprintf("查询菜单列表失败: %v", err)}, nil
	}

	list := make([]*apps.Menu, 0, len(menus))
	for _, m := range menus {
		list = append(list, menuToPb(m, tenantId))
	}

	return &apps.MenuListResp{
		Code:  200,
		Msg:   "success",
		List:  list,
		Total: int64(total),
	}, nil
}