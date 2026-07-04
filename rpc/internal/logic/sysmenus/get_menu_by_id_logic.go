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

type GetMenuByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMenuByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuByIdLogic {
	return &GetMenuByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMenuByIdLogic) GetMenuById(in *apps.IdReq) (*apps.MenuResp, error) {
	tenantId := mixins.GetCurrentTenantId(l.ctx)

	menu, err := l.svcCtx.DB.SysMenu.Query().
		Where(sysmenu.IDEQ(int64(in.Id))).
		Only(l.ctx)
	if err != nil {
		return &apps.MenuResp{Code: 500, Msg: fmt.Sprintf("获取菜单失败: %v", err)}, nil
	}

	return &apps.MenuResp{
		Code: 200,
		Msg:  "success",
		Data: menuToPb(menu, tenantId),
	}, nil
}

func menuToPb(m *ent.SysMenu, tenantId int64) *apps.Menu {
	sort := int32(m.Sort)
	menu := &apps.Menu{
		Id:        &m.ID,
		IdStr:     strPtr(fmt.Sprintf("%d", m.ID)),
		MenuType:  strPtr(string(m.MenuType)),
		Name:      &m.Name,
		ParentId:  &m.ParentID,
		ParentIdStr: strPtr(fmt.Sprintf("%d", m.ParentID)),
		Component: &m.Component,
		Path:      &m.Path,
		Icon:      &m.Icon,
		IsRedirect: &m.IsRedirect,
		Redirect:  &m.Redirect,
		Hidden:    &m.Hidden,
		Status:    strPtr(string(m.Status)),
		Sort:      &sort,
		Remark:    &m.Remark,
		TenantId:  &tenantId,
		TenantIdStr: strPtr(fmt.Sprintf("%d", tenantId)),
	}

	createdAt := m.CreatedAt.Unix()
	menu.CreatedAt = &createdAt
	menu.CreatedBy = &m.CreatedBy

	updatedAt := m.UpdatedAt.Unix()
	menu.UpdatedAt = &updatedAt
	menu.UpdatedBy = &m.UpdatedBy

	return menu
}

func strPtr(s string) *string {
	return &s
}