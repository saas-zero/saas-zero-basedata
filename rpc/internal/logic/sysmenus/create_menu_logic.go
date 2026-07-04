package sysmenuslogic

import (
	"context"
	"fmt"

	"github.com/saas-zero/saas-zero-basedata/ent/sysmenu"
	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMenuLogic {
	return &CreateMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateMenuLogic) CreateMenu(in *apps.MenuReq) (*apps.MenuResp, error) {
	if in.Name == nil || *in.Name == "" {
		return &apps.MenuResp{Code: 500, Msg: "菜单名称不能为空"}, nil
	}

	tenantId := mixins.GetCurrentTenantId(l.ctx)
	ctx := mixins.SetCurrentTenantId(l.ctx, tenantId)

	create := l.svcCtx.DB.SysMenu.Create().
		SetName(*in.Name)

	if in.MenuType != nil {
		create.SetMenuType(sysmenu.MenuType(*in.MenuType))
	}
	if in.ParentId != nil {
		create.SetParentID(*in.ParentId)
	}
	if in.Path != nil {
		create.SetPath(*in.Path)
	}
	if in.Component != nil {
		create.SetComponent(*in.Component)
	}
	if in.Icon != nil {
		create.SetIcon(*in.Icon)
	}
	if in.IsRedirect != nil {
		create.SetIsRedirect(*in.IsRedirect)
	}
	if in.Redirect != nil {
		create.SetRedirect(*in.Redirect)
	}
	if in.Hidden != nil {
		create.SetHidden(*in.Hidden)
	}
	if in.Status != nil {
		create.SetStatus(sysmenu.Status(*in.Status))
	}
	if in.Sort != nil {
		create.SetSort(uint32(*in.Sort))
	}
	if in.Remark != nil {
		create.SetRemark(*in.Remark)
	}

	menu, err := create.Save(ctx)
	if err != nil {
		return &apps.MenuResp{Code: 500, Msg: fmt.Sprintf("创建菜单失败: %v", err)}, nil
	}

	return &apps.MenuResp{
		Code: 200,
		Msg:  "创建成功",
		Data: menuToPb(menu, tenantId),
	}, nil
}