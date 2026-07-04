package sysmenuslogic

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/ent"
	"github.com/saas-zero/saas-zero-basedata/ent/sysmenu"
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

func (l *GetMenuTreeLogic) GetMenuTree(_ *apps.EmptyReq) (*apps.MenuTreeResp, error) {
	allMenus, err := l.svcCtx.DB.SysMenu.ActiveQuery().
		Order(ent.Asc(sysmenu.FieldSort)).
		All(l.ctx)
	if err != nil {
		return nil, err
	}

	tree := buildMenuTree(allMenus, 0)
	return &apps.MenuTreeResp{
		Code: 200,
		Msg:  "success",
		Data: tree,
	}, nil
}
