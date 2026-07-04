package systenantslogic

import (
	"context"
	"time"

	"github.com/saas-zero/saas-zero-basedata/ent/systenant"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"
	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
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
	userId := mixins.GetCurrentUserId(l.ctx)
	userName := mixins.GetCurrentUserName(l.ctx)
	ctx := mixins.SetCurrentUserId(l.ctx, userId)
	ctx = mixins.SetCurrentUserName(ctx, userName)

	update := l.svcCtx.DB.SysTenant.UpdateOneID(in.GetId())
	if in.Name != nil {
		update.SetName(in.GetName())
	}
	if in.Code != nil {
		update.SetCode(in.GetCode())
	}
	if in.AdminId != nil {
		update.SetAdminID(in.GetAdminId())
	}
	if in.ParentId != nil && in.GetParentId() > 0 {
		update.SetParentID(in.GetParentId())
	}
	if in.PackageId != nil && in.GetPackageId() > 0 {
		update.SetPackageID(in.GetPackageId())
	}
	if in.ExpiredAt != nil {
		update.SetExpiredAt(time.UnixMilli(in.GetExpiredAt()))
	}
	if in.Status != nil {
		update.SetStatus(systenant.Status(in.GetStatus()))
	}
	if in.Remark != nil {
		update.SetRemark(in.GetRemark())
	}

	result, err := update.Save(ctx)
	if err != nil {
		return nil, err
	}

	t, _ := l.svcCtx.DB.SysTenant.Query().Where(systenant.IDEQ(result.ID)).WithSysPackage().Only(ctx)
	return &apps.TenantResp{
		Code: 200,
		Msg:  "success",
		Data: tenantToResp(t),
	}, nil
}
