package sysmenuslogic

import (
	"context"
	"time"

	"github.com/saas-zero/saas-zero-basedata/ent/sysmenu"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"
	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMenuLogic {
	return &DeleteMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteMenuLogic) DeleteMenu(in *apps.IdsReq) (*apps.EmptyResp, error) {
	userId := mixins.GetCurrentUserId(l.ctx)
	userName := mixins.GetCurrentUserName(l.ctx)
	ctx := mixins.SetCurrentUserId(l.ctx, userId)
	ctx = mixins.SetCurrentUserName(ctx, userName)

	for _, id := range in.GetIds() {
		childCount, err := l.svcCtx.DB.SysMenu.Query().Where(sysmenu.ParentIDEQ(id), sysmenu.DeletedAtIsNil()).Count(ctx)
		if err != nil {
			return nil, err
		}
		if childCount > 0 {
			return &apps.EmptyResp{Code: 400, Msg: "该菜单下存在子菜单，无法删除"}, nil
		}
	}

	_, err := l.svcCtx.DB.SysMenu.Update().
		Where(sysmenu.IDIn(in.GetIds()...)).
		SetDeletedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return &apps.EmptyResp{Code: 200, Msg: "success"}, nil
}
