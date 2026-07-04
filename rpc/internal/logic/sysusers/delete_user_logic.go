package sysuserslogic

import (
	"context"
	"fmt"
	"time"

	"github.com/saas-zero/saas-zero-basedata/ent/sysuser"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteUserLogic) DeleteUser(in *apps.IdsReq) (*apps.EmptyResp, error) {
	if len(in.Ids) == 0 {
		return &apps.EmptyResp{Code: 500, Msg: "用户ID不能为空"}, nil
	}

	_, err := l.svcCtx.DB.SysUser.Update().
		Where(sysuser.IDIn(in.Ids...)).
		SetDeletedAt(time.Now()).
		Save(l.ctx)
	if err != nil {
		return &apps.EmptyResp{Code: 500, Msg: fmt.Sprintf("删除用户失败: %v", err)}, nil
	}

	return &apps.EmptyResp{Code: 200, Msg: "删除成功"}, nil
}