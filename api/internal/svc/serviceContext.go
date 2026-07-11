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
	"github.com/saas-zero/saas-zero-common/pkg/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config       config.Config
	Redis        *redis.Client
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
	conn := zrpc.MustNewClient(c.Basedata, zrpc.WithUnaryClientInterceptor(authClientInterceptor))

	// Casbin enforcer initialization with graceful degradation.
	// If PostgreSQL or Casbin initialization fails, enforcer is nil and the
	// CasbinAuth middleware will allow all requests (fail-open).
	var enf *casbinapi.SyncedEnforcer
	db, err := sql.Open("postgres", c.CasbinPostgres.DataSource)
	if err != nil {
		log.Printf("warning: failed to open casbin db: %v (casbin disabled)", err)
	} else {
		enf, err = commcasbin.NewEnforcer(db, "casbin_rule")
		if err != nil {
			log.Printf("warning: failed to init casbin enforcer: %v (casbin disabled)", err)
		}
	}
	if enf != nil {
		// Background goroutine: periodically reload Casbin policies from DB.
		// Policies are updated by basedata-rpc's AssignApis RPC.
		// 30s interval for faster policy propagation during development.
		go func() {
			ticker := time.NewTicker(30 * time.Second)
			defer ticker.Stop()
			for range ticker.C {
				if err := enf.LoadPolicy(); err != nil {
					log.Printf("casbin reload policy error: %v", err)
				}
			}
		}()
	}

	// Redis client initialization with graceful degradation.
	// If Redis is unavailable, JWT validation (token existence + version check)
	// will be skipped, falling back to JWT signature verification only.
	var rds *redis.Client
	rds, err = redis.NewClient(c.Redis)
	if err != nil {
		log.Printf("warning: failed to init redis: %v (redis disabled)", err)
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
