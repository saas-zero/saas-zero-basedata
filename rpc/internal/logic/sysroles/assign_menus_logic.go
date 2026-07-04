package sysroleslogic

import (
	"context"
	"fmt"

	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AssignMenusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAssignMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignMenusLogic {
	return &AssignMenusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AssignMenusLogic) AssignMenus(in *apps.RoleReq) (*apps.EmptyResp, error) {
	if in.Id == nil || *in.Id <= 0 {
		return &apps.EmptyResp{Code: 500, Msg: "角色ID不能为空"}, nil
	}

	tx, err := l.svcCtx.DB.Tx(l.ctx)
	if err != nil {
		return &apps.EmptyResp{Code: 500, Msg: fmt.Sprintf("开启事务失败: %v", err)}, nil
	}

	_, err = tx.SysRole.UpdateOneID(int64(*in.Id)).ClearMenus().Save(l.ctx)
	if err != nil {
		tx.Rollback()
		return &apps.EmptyResp{Code: 500, Msg: fmt.Sprintf("清空菜单权限失败: %v", err)}, nil
	}

	if len(in.MenuIds) > 0 {
		_, err = tx.SysRole.UpdateOneID(int64(*in.Id)).AddMenuIDs(in.MenuIds...).Save(l.ctx)
		if err != nil {
			tx.Rollback()
			return &apps.EmptyResp{Code: 500, Msg: fmt.Sprintf("分配菜单权限失败: %v", err)}, nil
		}
	}

	if err := tx.Commit(); err != nil {
		return &apps.EmptyResp{Code: 500, Msg: fmt.Sprintf("提交事务失败: %v", err)}, nil
	}

	return &apps.EmptyResp{Code: 200, Msg: "菜单权限分配成功"}, nil
}