package systenantslogic

import (
	"context"
	"fmt"

	"github.com/saas-zero/saas-zero-basedata/ent/systenant"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangeTenantStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChangeTenantStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeTenantStatusLogic {
	return &ChangeTenantStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ChangeTenantStatusLogic) ChangeTenantStatus(in *apps.TenantReq) (*apps.EmptyResp, error) {
	if in.Id == nil || *in.Id <= 0 {
		return &apps.EmptyResp{Code: 500, Msg: "租户ID不能为空"}, nil
	}

	if in.Status == nil || *in.Status == "" {
		return &apps.EmptyResp{Code: 500, Msg: "状态不能为空"}, nil
	}

	_, err := l.svcCtx.DB.SysTenant.UpdateOneID(int64(*in.Id)).
		SetStatus(systenant.Status(*in.Status)).
		Save(l.ctx)
	if err != nil {
		return &apps.EmptyResp{Code: 500, Msg: fmt.Sprintf("修改租户状态失败: %v", err)}, nil
	}

	return &apps.EmptyResp{Code: 200, Msg: "状态修改成功"}, nil
}