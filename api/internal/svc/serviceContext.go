package svc

import (
	"github.com/saas-zero/saas-zero-basedata/api/internal/config"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config       config.Config
	SysUsers     apps.SysUsersClient
	SysRoles     apps.SysRolesClient
	SysDepts     apps.SysDeptsClient
	SysMenus     apps.SysMenusClient
	SysDicts     apps.SysDictsClient
	SysDictDatas apps.SysDictDatasClient
	SysTenants   apps.SysTenantsClient
	SysPackages  apps.SysPackagesClient
	SysApis      apps.SysApisClient
	SysLogs      apps.SysLogsClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := zrpc.MustNewClient(c.Basedata)
	return &ServiceContext{
		Config:       c,
		SysUsers:     apps.NewSysUsersClient(conn.Conn()),
		SysRoles:     apps.NewSysRolesClient(conn.Conn()),
		SysDepts:     apps.NewSysDeptsClient(conn.Conn()),
		SysMenus:     apps.NewSysMenusClient(conn.Conn()),
		SysDicts:     apps.NewSysDictsClient(conn.Conn()),
		SysDictDatas: apps.NewSysDictDatasClient(conn.Conn()),
		SysTenants:   apps.NewSysTenantsClient(conn.Conn()),
		SysPackages:  apps.NewSysPackagesClient(conn.Conn()),
		SysApis:      apps.NewSysApisClient(conn.Conn()),
		SysLogs:      apps.NewSysLogsClient(conn.Conn()),
	}
}
