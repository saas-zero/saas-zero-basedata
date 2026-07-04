package sysdeptslogic

import (
	"context"
	"time"

	"github.com/saas-zero/saas-zero-basedata/ent/sysdept"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"
	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
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
	userId := mixins.GetCurrentUserId(l.ctx)
	userName := mixins.GetCurrentUserName(l.ctx)
	ctx := mixins.SetCurrentUserId(l.ctx, userId)
	ctx = mixins.SetCurrentUserName(ctx, userName)

	for _, id := range in.GetIds() {
		childCount, err := l.svcCtx.DB.SysDept.Query().Where(sysdept.ParentIDEQ(id), sysdept.DeletedAtIsNil()).Count(ctx)
		if err != nil {
			return nil, err
		}
		if childCount > 0 {
			return &apps.EmptyResp{Code: 400, Msg: "该部门下存在子部门，无法删除"}, nil
		}
	}

	_, err := l.svcCtx.DB.SysDept.Update().
		Where(sysdept.IDIn(in.GetIds()...)).
		SetDeletedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return &apps.EmptyResp{Code: 200, Msg: "success"}, nil
}
