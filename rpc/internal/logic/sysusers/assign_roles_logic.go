package sysuserslogic

import (
	"context"
	"fmt"

	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AssignRolesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAssignRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignRolesLogic {
	return &AssignRolesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AssignRolesLogic) AssignRoles(in *apps.UserReq) (*apps.EmptyResp, error) {
	if in.Id == nil || *in.Id <= 0 {
		return &apps.EmptyResp{Code: 500, Msg: "用户ID不能为空"}, nil
	}

	tx, err := l.svcCtx.DB.Tx(l.ctx)
	if err != nil {
		return &apps.EmptyResp{Code: 500, Msg: fmt.Sprintf("开启事务失败: %v", err)}, nil
	}

	_, err = tx.SysUser.UpdateOneID(int64(*in.Id)).ClearRoles().Save(l.ctx)
	if err != nil {
		tx.Rollback()
		return &apps.EmptyResp{Code: 500, Msg: fmt.Sprintf("清空角色失败: %v", err)}, nil
	}

	if len(in.RoleIds) > 0 {
		_, err = tx.SysUser.UpdateOneID(int64(*in.Id)).AddRoleIDs(in.RoleIds...).Save(l.ctx)
		if err != nil {
			tx.Rollback()
			return &apps.EmptyResp{Code: 500, Msg: fmt.Sprintf("分配角色失败: %v", err)}, nil
		}
	}

	if err := tx.Commit(); err != nil {
		return &apps.EmptyResp{Code: 500, Msg: fmt.Sprintf("提交事务失败: %v", err)}, nil
	}

	return &apps.EmptyResp{Code: 200, Msg: "角色分配成功"}, nil
}