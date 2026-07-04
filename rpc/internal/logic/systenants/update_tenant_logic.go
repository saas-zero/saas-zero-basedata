package systenantslogic

import (
	"context"
	"fmt"
	"time"

	"github.com/saas-zero/saas-zero-basedata/ent/systenant"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTenantLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateTenantLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTenantLogic {
	return &UpdateTenantLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateTenantLogic) UpdateTenant(in *apps.TenantReq) (*apps.TenantResp, error) {
	if in.Id == nil || *in.Id <= 0 {
		return &apps.TenantResp{Code: 500, Msg: "租户ID不能为空"}, nil
	}

	update := l.svcCtx.DB.SysTenant.UpdateOneID(int64(*in.Id))

	if in.Name != nil {
		update.SetName(*in.Name)
	}
	if in.Code != nil {
		update.SetCode(*in.Code)
	}
	if in.AdminId != nil {
		update.SetAdminID(*in.AdminId)
	}
	if in.ParentId != nil {
		update.SetParentID(*in.ParentId)
	}
	if in.PackageId != nil {
		update.SetPackageID(*in.PackageId)
	}
	if in.ExpiredAt != nil {
		update.SetExpiredAt(time.Unix(*in.ExpiredAt, 0))
	}
	if in.Status != nil {
		update.SetStatus(systenant.Status(*in.Status))
	}
	if in.Remark != nil {
		update.SetRemark(*in.Remark)
	}

	tenant, err := update.Save(l.ctx)
	if err != nil {
		return &apps.TenantResp{Code: 500, Msg: fmt.Sprintf("更新租户失败: %v", err)}, nil
	}

	tenant, err = l.svcCtx.DB.SysTenant.Query().
		Where(systenant.IDEQ(tenant.ID)).
		WithSysPackage().
		Only(l.ctx)
	if err != nil {
		return &apps.TenantResp{Code: 500, Msg: fmt.Sprintf("获取租户失败: %v", err)}, nil
	}

	return &apps.TenantResp{
		Code: 200,
		Msg:  "更新成功",
		Data: tenantToPb(tenant),
	}, nil
}