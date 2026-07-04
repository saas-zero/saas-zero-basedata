package sysmenuslogic

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/ent"
	"github.com/saas-zero/saas-zero-basedata/ent/sysmenu"
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
	query := l.svcCtx.DB.SysMenu.ActiveQuery()
	if in.GetName() != "" {
		query = query.Where(sysmenu.NameContains(in.GetName()))
	}
	if in.GetStatus() != "" {
		query = query.Where(sysmenu.StatusEQ(sysmenu.Status(in.GetStatus())))
	}

	total, err := query.Count(l.ctx)
	if err != nil {
		return nil, err
	}

	page := int(in.GetPage())
	size := int(in.GetSize())
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	menus, err := query.
		Offset((page - 1) * size).
		Limit(size).
		Order(ent.Asc(sysmenu.FieldSort)).
		All(l.ctx)
	if err != nil {
		return nil, err
	}

	list := make([]*apps.Menu, len(menus))
	for i, m := range menus {
		list[i] = menuToResp(m)
	}
	return &apps.MenuListResp{
		Code:  200,
		Msg:   "success",
		List:  list,
		Total: int64(total),
	}, nil
}
