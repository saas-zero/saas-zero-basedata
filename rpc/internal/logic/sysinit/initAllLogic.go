package sysinitlogic

import (
	"context"
	"strconv"

	"github.com/saas-zero/saas-zero-basedata/ent/sysapi"
	"github.com/saas-zero/saas-zero-basedata/ent/sysdept"
	"github.com/saas-zero/saas-zero-basedata/ent/sysmenu"
	"github.com/saas-zero/saas-zero-basedata/ent/syspackage"
	"github.com/saas-zero/saas-zero-basedata/ent/sysrole"
	"github.com/saas-zero/saas-zero-basedata/ent/systenant"
	"github.com/saas-zero/saas-zero-basedata/ent/sysuser"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"
	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type InitAllLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewInitAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InitAllLogic {
	return &InitAllLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

type seedApi struct {
	name string
	path string
}

type seedMenu struct {
	menuType  string
	name      string
	parentId  int64
	component string
	path      string
	icon      string
	sort      uint32
}

func (l *InitAllLogic) InitAll(_ *apps.EmptyReq) (*apps.EmptyResp, error) {
	tenantId := mixins.GetCurrentTenantId(l.ctx)
	userId := mixins.GetCurrentUserId(l.ctx)
	userName := mixins.GetCurrentUserName(l.ctx)

	ctx := mixins.SetCurrentTenantId(l.ctx, tenantId)
	ctx = mixins.SetCurrentUserId(ctx, userId)
	ctx = mixins.SetCurrentUserName(ctx, userName)

	seedApis := []seedApi{
		{name: "用户管理", path: "/system/user/*"},
		{name: "角色管理", path: "/system/role/*"},
		{name: "菜单管理", path: "/system/menu/*"},
		{name: "部门管理", path: "/system/dept/*"},
		{name: "字典管理", path: "/system/dict/*"},
		{name: "字典数据管理", path: "/system/dictData/*"},
		{name: "租户管理", path: "/system/tenant/*"},
		{name: "套餐管理", path: "/system/package/*"},
		{name: "API管理", path: "/system/api/*"},
	}

	seedMenus := []seedMenu{
		{menuType: "directory", name: "系统管理", path: "/system", icon: "Setting", sort: 1},
		{menuType: "menu", name: "用户管理", parentId: 1, component: "system/user/index", path: "/system/user", icon: "User", sort: 1},
		{menuType: "menu", name: "角色管理", parentId: 1, component: "system/role/index", path: "/system/role", icon: "SafetyCertificate", sort: 2},
		{menuType: "menu", name: "菜单管理", parentId: 1, component: "system/menu/index", path: "/system/menu", icon: "Menu", sort: 3},
		{menuType: "menu", name: "部门管理", parentId: 1, component: "system/dept/index", path: "/system/dept", icon: "Apartment", sort: 4},
		{menuType: "menu", name: "字典管理", parentId: 1, component: "system/dict/index", path: "/system/dict", icon: "Book", sort: 5},
		{menuType: "menu", name: "租户管理", parentId: 1, component: "system/tenant/index", path: "/system/tenant", icon: "Team", sort: 6},
		{menuType: "menu", name: "套餐管理", parentId: 1, component: "system/package/index", path: "/system/package", icon: "Dollar", sort: 7},
	}

	tx, err := l.svcCtx.DB.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// 1. Create APIs
	apiIds := make([]int64, 0, len(seedApis))
	apiPaths := make([]string, 0, len(seedApis))
	for _, a := range seedApis {
		api, err := tx.SysApi.Create().
			SetStatus(sysapi.StatusActive).
			SetAPIName(a.name).
			SetAPIType(sysapi.APITypeGroup).
			SetAPIPath(a.path).
			Save(ctx)
		if err != nil {
			return nil, err
		}
		apiIds = append(apiIds, api.ID)
		apiPaths = append(apiPaths, a.path)
	}

	// 2. Create Menus
	menuIds := make([]int64, 0, len(seedMenus))
	for _, m := range seedMenus {
		create := tx.SysMenu.Create().
			SetStatus(sysmenu.StatusActive).
			SetSort(m.sort).
			SetMenuType(sysmenu.MenuType(m.menuType)).
			SetName(m.name).
			SetPath(m.path).
			SetIcon(m.icon)
		if m.component != "" {
			create.SetComponent(m.component)
		}
		if m.parentId > 0 {
			create.SetParentID(m.parentId)
		}
		menu, err := create.Save(ctx)
		if err != nil {
			return nil, err
		}
		menuIds = append(menuIds, menu.ID)
	}

	// 3. Create Package (linked to all menus and APIs)
	pkg, err := tx.SysPackage.Create().
		SetName("标准版").
		SetCode("standard").
		SetSort(1).
		SetStatus(syspackage.StatusActive).
		AddMenuIDs(menuIds...).
		AddAPIIDs(apiIds...).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	// 4. Create Tenant
	tenant, err := tx.SysTenant.Create().
		SetName("默认租户").
		SetCode("default").
		SetAdminID(1).
		SetPackageID(pkg.ID).
		SetStatus(systenant.StatusActive).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	// Update context with actual tenant ID so subsequent TenantMixin hooks use the correct value
	ctx = mixins.SetCurrentTenantId(ctx, tenant.ID)

	// 5. Create Role (linked to all menus)
	role, err := tx.SysRole.Create().
		SetName("超级管理员").
		SetCode("admin").
		SetSort(1).
		SetStatus(sysrole.StatusActive).
		AddMenuIDs(menuIds...).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	// 6. Create Default Department
	dept, err := tx.SysDept.Create().
		SetName("默认部门").
		SetSort(1).
		SetStatus(sysdept.StatusActive).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	// 7. Create User (linked to role and dept)
	hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	_, err = tx.SysUser.Create().
		SetUsername("admin").
		SetPassword(string(hash)).
		SetNickname("系统管理员").
		SetDeptID(dept.ID).
		SetStatus(sysuser.StatusActive).
		AddRoleIDs(role.ID).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// 8. Casbin policies (after ent tx commits)
	dom := strconv.FormatInt(tenant.ID, 10)
	for i, id := range apiIds {
		l.svcCtx.Enforcer.AddPolicy("admin", dom, apiPaths[i], ".*", strconv.FormatInt(id, 10))
	}

	return &apps.EmptyResp{Code: 200, Msg: "success"}, nil
}
