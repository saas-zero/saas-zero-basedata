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

type UpdateMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMenuLogic {
	return &UpdateMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateMenuLogic) UpdateMenu(in *apps.MenuReq) (*apps.MenuResp, error) {
	if in.Id == nil || *in.Id <= 0 {
		return &apps.MenuResp{Code: 500, Msg: "菜单ID不能为空"}, nil
	}

	tenantId := mixins.GetCurrentTenantId(l.ctx)

	update := l.svcCtx.DB.SysMenu.UpdateOneID(int64(*in.Id))

	if in.Name != nil {
		update.SetName(*in.Name)
	}
	if in.MenuType != nil {
		update.SetMenuType(sysmenu.MenuType(*in.MenuType))
	}
	if in.ParentId != nil {
		update.SetParentID(*in.ParentId)
	}
	if in.Path != nil {
		update.SetPath(*in.Path)
	}
	if in.Component != nil {
		update.SetComponent(*in.Component)
	}
	if in.Icon != nil {
		update.SetIcon(*in.Icon)
	}
	if in.IsRedirect != nil {
		update.SetIsRedirect(*in.IsRedirect)
	}
	if in.Redirect != nil {
		update.SetRedirect(*in.Redirect)
	}
	if in.Hidden != nil {
		update.SetHidden(*in.Hidden)
	}
	if in.Status != nil {
		update.SetStatus(sysmenu.Status(*in.Status))
	}
	if in.Sort != nil {
		update.SetSort(uint32(*in.Sort))
	}
	if in.Remark != nil {
		update.SetRemark(*in.Remark)
	}

	menu, err := update.Save(l.ctx)
	if err != nil {
		return &apps.MenuResp{Code: 500, Msg: fmt.Sprintf("更新菜单失败: %v", err)}, nil
	}

	return &apps.MenuResp{
		Code: 200,
		Msg:  "更新成功",
		Data: menuToPb(menu, tenantId),
	}, nil
}