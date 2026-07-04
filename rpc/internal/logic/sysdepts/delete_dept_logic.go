package sysdeptslogic

import (
	"context"
	"fmt"
	"time"

	"github.com/saas-zero/saas-zero-basedata/ent/sysdept"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteDeptLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteDeptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDeptLogic {
	return &DeleteDeptLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteDeptLogic) DeleteDept(in *apps.IdsReq) (*apps.EmptyResp, error) {
	if len(in.Ids) == 0 {
		return &apps.EmptyResp{Code: 500, Msg: "部门ID不能为空"}, nil
	}

	_, err := l.svcCtx.DB.SysDept.Update().
		Where(sysdept.IDIn(in.Ids...)).
		SetDeletedAt(time.Now()).
		Save(l.ctx)
	if err != nil {
		return &apps.EmptyResp{Code: 500, Msg: fmt.Sprintf("删除部门失败: %v", err)}, nil
	}

	return &apps.EmptyResp{Code: 200, Msg: "删除成功"}, nil
}