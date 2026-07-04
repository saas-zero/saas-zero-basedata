package sysroleslogic

import (
	"context"
	"fmt"

	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AssignApisLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAssignApisLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignApisLogic {
	return &AssignApisLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AssignApisLogic) AssignApis(in *apps.RoleReq) (*apps.EmptyResp, error) {
	if in.Id == nil || *in.Id <= 0 {
		return &apps.EmptyResp{Code: 500, Msg: "角色ID不能为空"}, nil
	}

	tx, err := l.svcCtx.DB.Tx(l.ctx)
	if err != nil {
		return &apps.EmptyResp{Code: 500, Msg: fmt.Sprintf("开启事务失败: %v", err)}, nil
	}

	_, err = tx.SysRole.UpdateOneID(int64(*in.Id)).ClearApis().Save(l.ctx)
	if err != nil {
		tx.Rollback()
		return &apps.EmptyResp{Code: 500, Msg: fmt.Sprintf("清空API权限失败: %v", err)}, nil
	}

	if len(in.ApiIds) > 0 {
		_, err = tx.SysRole.UpdateOneID(int64(*in.Id)).AddAPIIDs(in.ApiIds...).Save(l.ctx)
		if err != nil {
			tx.Rollback()
			return &apps.EmptyResp{Code: 500, Msg: fmt.Sprintf("分配API权限失败: %v", err)}, nil
		}
	}

	if err := tx.Commit(); err != nil {
		return &apps.EmptyResp{Code: 500, Msg: fmt.Sprintf("提交事务失败: %v", err)}, nil
	}

	return &apps.EmptyResp{Code: 200, Msg: "API权限分配成功"}, nil
}