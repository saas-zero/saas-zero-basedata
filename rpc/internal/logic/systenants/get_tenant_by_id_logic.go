package systenantslogic

import (
	"context"
	"fmt"

	"github.com/saas-zero/saas-zero-basedata/ent"
	"github.com/saas-zero/saas-zero-basedata/ent/systenant"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTenantByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTenantByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTenantByIdLogic {
	return &GetTenantByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTenantByIdLogic) GetTenantById(in *apps.IdReq) (*apps.TenantResp, error) {
	tenant, err := l.svcCtx.DB.SysTenant.Query().
		Where(systenant.IDEQ(int64(in.Id))).
		WithSysPackage().
		Only(l.ctx)
	if err != nil {
		return &apps.TenantResp{Code: 500, Msg: fmt.Sprintf("获取租户失败: %v", err)}, nil
	}

	return &apps.TenantResp{
		Code: 200,
		Msg:  "success",
		Data: tenantToPb(tenant),
	}, nil
}

func tenantToPb(t *ent.SysTenant) *apps.Tenant {
	tenant := &apps.Tenant{
		Id:       &t.ID,
		IdStr:    strPtr(fmt.Sprintf("%d", t.ID)),
		Name:     &t.Name,
		Code:     &t.Code,
		Status:   strPtr(string(t.Status)),
		Remark:   &t.Remark,
	}

	if t.PackageID != 0 {
		tenant.PackageId = &t.PackageID
		tenant.PackageIdStr = strPtr(fmt.Sprintf("%d", t.PackageID))
	}

	if t.Edges.SysPackage != nil {
		tenant.PackageName = &t.Edges.SysPackage.Name
	}

	if !t.ExpiredAt.IsZero() {
		expiredAt := t.ExpiredAt.Unix()
		tenant.ExpiredAt = &expiredAt
	}

	createdAt := t.CreatedAt.Unix()
	tenant.CreatedAt = &createdAt
	tenant.CreatedBy = &t.CreatedBy

	updatedAt := t.UpdatedAt.Unix()
	tenant.UpdatedAt = &updatedAt
	tenant.UpdatedBy = &t.UpdatedBy

	return tenant
}

func strPtr(s string) *string {
	return &s
}