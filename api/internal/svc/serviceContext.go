package svc

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"

	casbinapi "github.com/casbin/casbin/v2"
	"github.com/saas-zero/saas-zero-basedata/api/internal/config"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	commcasbin "github.com/saas-zero/saas-zero-common/pkg/casbin"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config       config.Config
	Redis        *redis.Redis
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
	SysInit      apps.SysInitClient
	Enforcer     *casbinapi.SyncedEnforcer
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := zrpc.MustNewClient(c.Basedata)

	db, err := sql.Open("postgres", c.CasbinPostgres.DataSource)
	if err != nil {
		log.Fatalf("failed to open casbin db: %v", err)
	}
	enf, err := commcasbin.NewEnforcer(db, "casbin_rule")
	if err != nil {
		log.Fatalf("failed to init casbin enforcer: %v", err)
	}
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			if err := enf.LoadPolicy(); err != nil {
				log.Printf("casbin reload policy error: %v", err)
			}
		}
	}()

	rds, err := redis.NewRedis(c.Redis)
	if err != nil {
		log.Fatalf("failed to init redis: %v", err)
	}

	return &ServiceContext{
		Config:       c,
		Redis:        rds,
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
		SysInit:      apps.NewSysInitClient(conn.Conn()),
		Enforcer:     enf,
	}
}
