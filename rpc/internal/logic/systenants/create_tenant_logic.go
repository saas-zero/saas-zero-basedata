package systenantslogic

import (
	"context"
	"fmt"
	"time"

	"github.com/saas-zero/saas-zero-basedata/ent"
	"github.com/saas-zero/saas-zero-basedata/ent/syspackage"
	"github.com/saas-zero/saas-zero-basedata/ent/sysapi"
	"github.com/saas-zero/saas-zero-basedata/ent/sysmenu"
	"github.com/saas-zero/saas-zero-basedata/ent/systenant"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateTenantLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateTenantLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTenantLogic {
	return &CreateTenantLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateTenantLogic) CreateTenant(in *apps.TenantReq) (*apps.TenantResp, error) {
	if in.Name == nil || *in.Name == "" {
		return &apps.TenantResp{Code: 500, Msg: "租户名称不能为空"}, nil
	}

	if in.Code == nil || *in.Code == "" {
		return &apps.TenantResp{Code: 500, Msg: "租户编码不能为空"}, nil
	}

	if in.AdminId == nil || *in.AdminId <= 0 {
		return &apps.TenantResp{Code: 500, Msg: "管理员ID不能为空"}, nil
	}

	tx, err := l.svcCtx.DB.Tx(l.ctx)
	if err != nil {
		return &apps.TenantResp{Code: 500, Msg: fmt.Sprintf("开启事务失败: %v", err)}, nil
	}

	create := tx.SysTenant.Create().
		SetName(*in.Name).
		SetCode(*in.Code).
		SetAdminID(*in.AdminId)

	if in.ParentId != nil {
		create.SetParentID(*in.ParentId)
	}
	if in.PackageId != nil {
		create.SetPackageID(*in.PackageId)
	}
	if in.ExpiredAt != nil {
		create.SetExpiredAt(time.Unix(*in.ExpiredAt, 0))
	}
	if in.Status != nil {
		create.SetStatus(systenant.Status(*in.Status))
	}
	if in.Remark != nil {
		create.SetRemark(*in.Remark)
	}

	tenant, err := create.Save(l.ctx)
	if err != nil {
		tx.Rollback()
		return &apps.TenantResp{Code: 500, Msg: fmt.Sprintf("创建租户失败: %v", err)}, nil
	}

	if err := l.createDefaultAdminRole(tx, tenant.ID, *in.PackageId, *in.Name); err != nil {
		tx.Rollback()
		return &apps.TenantResp{Code: 500, Msg: fmt.Sprintf("创建默认角色失败: %v", err)}, nil
	}

	if err := tx.Commit(); err != nil {
		return &apps.TenantResp{Code: 500, Msg: fmt.Sprintf("提交事务失败: %v", err)}, nil
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
		Msg:  "创建成功",
		Data: tenantToPb(tenant),
	}, nil
}

func (l *CreateTenantLogic) createDefaultAdminRole(tx *ent.Tx, tenantId, packageId int64, tenantName string) error {
	if packageId <= 0 {
		return nil
	}

	menuIds, err := tx.SysMenu.Query().
		Where(sysmenu.HasPackagesWith(syspackage.IDEQ(packageId))).
		Select(sysmenu.FieldID).
		All(l.ctx)
	if err != nil {
		return err
	}

	apiIds, err := tx.SysApi.Query().
		Where(sysapi.HasPackagesWith(syspackage.IDEQ(packageId))).
		Select(sysapi.FieldID).
		All(l.ctx)
	if err != nil {
		return err
	}

	menuIDList := make([]int64, 0, len(menuIds))
	for _, m := range menuIds {
		menuIDList = append(menuIDList, m.ID)
	}

	apiIDList := make([]int64, 0, len(apiIds))
	for _, a := range apiIds {
		apiIDList = append(apiIDList, a.ID)
	}

	_, err = tx.SysRole.Create().
		SetName(fmt.Sprintf("%s管理员", tenantName)).
		SetCode(fmt.Sprintf("tenant_%d_admin", tenantId)).
		SetIsSystem(true).
		AddMenuIDs(menuIDList...).
		AddAPIIDs(apiIDList...).
		Save(l.ctx)
	if err != nil {
		return err
	}

	return nil
}