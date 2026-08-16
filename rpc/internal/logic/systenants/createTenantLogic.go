package systenantslogic

import (
	"context"
	"strings"

	"github.com/saas-zero/saas-zero-basedata/ent/syspackage"
	"github.com/saas-zero/saas-zero-basedata/ent/sysrole"
	"github.com/saas-zero/saas-zero-basedata/ent/systenant"
	"github.com/saas-zero/saas-zero-basedata/ent/sysuser"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"
	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
	"github.com/saas-zero/saas-zero-common/pkg/errno"
	idutil "github.com/saas-zero/saas-zero-common/pkg/id"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/proto"
	"time"
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

// CreateTenant 创建租户并完成开通闭环：
//  1. 创建租户（关联套餐）
//  2. 创建默认角色（code=admin），继承套餐的菜单
//  3. 创建管理员用户（表单传入账号密码），分配默认角色
//  4. 回填租户 admin_id
//  5. 同步 Casbin：默认角色 + 新租户 dom 的 API 策略（来自套餐 API 关联）
func (l *CreateTenantLogic) CreateTenant(in *apps.TenantReq) (*apps.TenantResp, error) {
	if in.GetPackageId() <= 0 {
		return &apps.TenantResp{Code: int32(errno.InvalidParam.Code), Msg: "请选择套餐"}, nil
	}
	if in.GetUsername() == "" || in.GetPassword() == "" {
		return &apps.TenantResp{Code: int32(errno.InvalidParam.Code), Msg: "请填写管理员账号和密码"}, nil
	}

	userId := mixins.GetCurrentUserId(l.ctx)
	userName := mixins.GetCurrentUserName(l.ctx)
	ctx := mixins.SetCurrentUserId(l.ctx, userId)
	ctx = mixins.SetCurrentUserName(ctx, userName)

	tx, err := l.svcCtx.DB.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// 1. 创建租户（admin_id 先占位为操作者，创建管理员用户后回填）
	tenant, err := tx.SysTenant.Create().
		SetName(in.GetName()).
		SetCode(in.GetCode()).
		SetPackageID(in.GetPackageId()).
		SetStatus(systenant.Status(in.GetStatus())).
		SetAdminID(userId).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	if in.GetExpiredAt() > 0 {
		if err := tx.SysTenant.UpdateOneID(tenant.ID).SetExpiredAt(time.UnixMilli(in.GetExpiredAt())).Exec(ctx); err != nil {
			return nil, err
		}
	}
	if in.GetRemark() != "" {
		if err := tx.SysTenant.UpdateOneID(tenant.ID).SetRemark(in.GetRemark()).Exec(ctx); err != nil {
			return nil, err
		}
	}
	// 后续创建的角色/用户归属新租户
	ctx = mixins.SetCurrentTenantId(ctx, tenant.ID)

	// 2. 查套餐的菜单/API 关联
	pkg, err := tx.SysPackage.Query().
		Where(syspackage.IDEQ(in.GetPackageId())).
		WithMenus().WithApis().
		Only(ctx)
	if err != nil {
		return nil, err
	}
	menuIDs := make([]int64, 0, len(pkg.Edges.Menus))
	for _, m := range pkg.Edges.Menus {
		menuIDs = append(menuIDs, m.ID)
	}

	// 3. 创建默认角色（code=admin），继承套餐菜单
	role, err := tx.SysRole.Create().
		SetName("管理员").
		SetCode("admin").
		SetSort(1).
		SetStatus(sysrole.StatusActive).
		AddMenuIDs(menuIDs...).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	// 4. 创建管理员用户（bcrypt 密码），分配默认角色
	hash, err := bcrypt.GenerateFromPassword([]byte(in.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	admin, err := tx.SysUser.Create().
		SetUsername(in.GetUsername()).
		SetPassword(string(hash)).
		SetNickname(in.GetUsername()).
		SetStatus(sysuser.StatusActive).
		AddRoleIDs(role.ID).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	// 5. 回填租户 admin_id
	if err := tx.SysTenant.UpdateOneID(tenant.ID).SetAdminID(admin.ID).Exec(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// 6. 同步 Casbin：默认角色 + 新租户 dom 的 API 策略（套餐 API）
	dom := idutil.ToString(tenant.ID)
	if l.svcCtx.Enforcer != nil {
		for _, api := range pkg.Edges.Apis {
			if _, err := l.svcCtx.Enforcer.AddPolicy("admin", dom, api.APIPath,
				strings.ToUpper(string(api.APIMethod)), idutil.ToString(api.ID)); err != nil {
				logx.Errorf("createTenant: add casbin policy failed: %v", err)
			}
		}
	}

	return &apps.TenantResp{
		Code: int32(errno.Success.Code),
		Msg:  errno.Success.Msg,
		Data: &apps.Tenant{
			Id:    proto.Int64(tenant.ID),
			IdStr: proto.String(idutil.ToString(tenant.ID)),
		},
	}, nil
}
