package sysinitlogic

import (
	"context"
	"github.com/saas-zero/saas-zero-common/pkg/id"

	"github.com/saas-zero/saas-zero-basedata/ent"
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
	"github.com/saas-zero/saas-zero-common/pkg/errno"
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
		{name: "日志管理", path: "/system/log/*"},
	}

	seedMenus := []struct {
		menuType  string
		name      string
		parentIdx int
		component string
		path      string
		icon      string
		sort      uint32
	}{
		{menuType: "menu", name: "控制台", parentIdx: -1, component: "dashboard/index", path: "/dashboard", icon: "Dashboard", sort: 1},
		{menuType: "directory", name: "系统管理", parentIdx: -1, path: "/system", icon: "Setting", sort: 2},
		{menuType: "menu", name: "用户管理", parentIdx: 1, component: "system/user/index", path: "/system/user", icon: "User", sort: 1},
		{menuType: "menu", name: "角色管理", parentIdx: 1, component: "system/role/index", path: "/system/role", icon: "SafetyCertificate", sort: 2},
		{menuType: "menu", name: "菜单管理", parentIdx: 1, component: "system/menu/index", path: "/system/menu", icon: "Menu", sort: 3},
		{menuType: "menu", name: "部门管理", parentIdx: 1, component: "system/dept/index", path: "/system/dept", icon: "Apartment", sort: 4},
		{menuType: "menu", name: "字典管理", parentIdx: -1, component: "dict/index", path: "/dict", icon: "Book", sort: 3},
		{menuType: "directory", name: "租户管理", parentIdx: -1, path: "/tenant", icon: "Team", sort: 4},
		{menuType: "menu", name: "租户列表", parentIdx: 7, component: "tenant/list/index", path: "/tenant/list", icon: "Team", sort: 1},
		{menuType: "menu", name: "套餐管理", parentIdx: 7, component: "tenant/package/index", path: "/tenant/package", icon: "Dollar", sort: 2},
		{menuType: "menu", name: "API管理", parentIdx: -1, component: "api/index", path: "/api", icon: "Code", sort: 5},
		{menuType: "directory", name: "日志管理", parentIdx: -1, path: "/log", icon: "FileText", sort: 6},
		{menuType: "menu", name: "登录日志", parentIdx: 11, component: "log/login-log/index", path: "/log/login-log", icon: "Login", sort: 1},
		{menuType: "menu", name: "操作日志", parentIdx: 11, component: "log/operation-log/index", path: "/log/operation-log", icon: "SwapRight", sort: 2},
	}

	tx, err := l.svcCtx.DB.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// 1. Create or fetch APIs (keyed by unique path)
	apiPaths := make([]string, 0, len(seedApis))
	apiIds := make([]int64, 0, len(seedApis))
	for _, a := range seedApis {
		api, err := tx.SysApi.Query().Where(sysapi.APIPathEQ(a.path)).First(ctx)
		if ent.IsNotFound(err) {
			api, err = tx.SysApi.Create().
				SetStatus(sysapi.StatusActive).
				SetAPIName(a.name).
				SetAPIType(sysapi.APITypeGroup).
				SetAPIPath(a.path).
				Save(ctx)
			if err != nil {
				return nil, err
			}
		} else if err != nil {
			return nil, err
		}
		apiPaths = append(apiPaths, a.path)
		apiIds = append(apiIds, api.ID)
	}

	// 2. Create or fetch Menus (keyed by unique name)
	menuIdxToId := make(map[int]int64)
	menuIds := make([]int64, 0, len(seedMenus))
	for i, m := range seedMenus {
		menu, err := tx.SysMenu.Query().Where(sysmenu.NameEQ(m.name)).First(ctx)
		if ent.IsNotFound(err) {
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
			if m.parentIdx >= 0 {
				create.SetParentID(menuIdxToId[m.parentIdx])
			}
			menu, err = create.Save(ctx)
			if err != nil {
				return nil, err
			}
		} else if err != nil {
			return nil, err
		}
		menuIdxToId[i] = menu.ID
		menuIds = append(menuIds, menu.ID)
	}

	// 3. Create or fetch Package (code "standard")
	pkg, err := tx.SysPackage.Query().Where(syspackage.CodeEQ("standard")).First(ctx)
	if ent.IsNotFound(err) {
		pkg, err = tx.SysPackage.Create().
			SetName("标准套餐").
			SetCode("standard").
			SetSort(1).
			SetStatus(syspackage.StatusActive).
			AddMenuIDs(menuIds...).
			AddAPIIDs(apiIds...).
			Save(ctx)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	// 4. Create or fetch Tenant (code "default")
	tenant, err := tx.SysTenant.Query().Where(systenant.CodeEQ("default")).First(ctx)
	if ent.IsNotFound(err) {
		tenant, err = tx.SysTenant.Create().
			SetName("默认租户").
			SetCode("default").
			SetAdminID(1).
			SetPackageID(pkg.ID).
			SetStatus(systenant.StatusActive).
			Save(ctx)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	// Update context with actual tenant ID
	ctx = mixins.SetCurrentTenantId(ctx, tenant.ID)

	// 5. Create or fetch Role (code "admin")
	role, err := tx.SysRole.Query().Where(sysrole.CodeEQ("admin")).First(ctx)
	if ent.IsNotFound(err) {
		role, err = tx.SysRole.Create().
			SetName("超级管理员").
			SetCode("admin").
			SetSort(1).
			SetStatus(sysrole.StatusActive).
			AddMenuIDs(menuIds...).
			Save(ctx)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	// 6. Create or fetch Department (name "默认部门")
	dept, err := tx.SysDept.Query().Where(sysdept.NameEQ("默认部门")).First(ctx)
	if ent.IsNotFound(err) {
		dept, err = tx.SysDept.Create().
			SetName("默认部门").
			SetSort(1).
			SetStatus(sysdept.StatusActive).
			Save(ctx)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	// 7. Create or fetch User (username "admin")
	_, err = tx.SysUser.Query().Where(sysuser.UsernameEQ("admin")).First(ctx)
	if ent.IsNotFound(err) {
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
	} else if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// 8. Casbin policies: clear old ones for this role+tenant, then re-add
	dom := id.ToString(tenant.ID)
	if _, err := l.svcCtx.Enforcer.RemoveFilteredPolicy(0, "admin", dom); err != nil {
		logx.Errorf("initAll: failed to clear casbin policies: %v", err)
	}
	for i, apiId := range apiIds {
		if _, err := l.svcCtx.Enforcer.AddPolicy("admin", dom, apiPaths[i], ".*", id.ToString(apiId)); err != nil {
			logx.Errorf("initAll: failed to add casbin policy: %v", err)
		}
	}

	return &apps.EmptyResp{Code: int32(errno.Success.Code), Msg: errno.Success.Msg}, nil
}
