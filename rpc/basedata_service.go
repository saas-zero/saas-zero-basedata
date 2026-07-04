package main

import (
	"flag"
	"fmt"

	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/config"
	sysapisServer "github.com/saas-zero/saas-zero-basedata/rpc/internal/server/sysapis"
	sysdeptsServer "github.com/saas-zero/saas-zero-basedata/rpc/internal/server/sysdepts"
	sysdictdatasServer "github.com/saas-zero/saas-zero-basedata/rpc/internal/server/sysdictdatas"
	sysdictsServer "github.com/saas-zero/saas-zero-basedata/rpc/internal/server/sysdicts"
	syslogsServer "github.com/saas-zero/saas-zero-basedata/rpc/internal/server/syslogs"
	sysmenusServer "github.com/saas-zero/saas-zero-basedata/rpc/internal/server/sysmenus"
	syspackagesServer "github.com/saas-zero/saas-zero-basedata/rpc/internal/server/syspackages"
	sysrolesServer "github.com/saas-zero/saas-zero-basedata/rpc/internal/server/sysroles"
	systenantsServer "github.com/saas-zero/saas-zero-basedata/rpc/internal/server/systenants"
	sysusersServer "github.com/saas-zero/saas-zero-basedata/rpc/internal/server/sysusers"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/basedata_service.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		apps.RegisterSysUsersServer(grpcServer, sysusersServer.NewSysUsersServer(ctx))
		apps.RegisterSysRolesServer(grpcServer, sysrolesServer.NewSysRolesServer(ctx))
		apps.RegisterSysDeptsServer(grpcServer, sysdeptsServer.NewSysDeptsServer(ctx))
		apps.RegisterSysMenusServer(grpcServer, sysmenusServer.NewSysMenusServer(ctx))
		apps.RegisterSysDictsServer(grpcServer, sysdictsServer.NewSysDictsServer(ctx))
		apps.RegisterSysDictDatasServer(grpcServer, sysdictdatasServer.NewSysDictDatasServer(ctx))
		apps.RegisterSysTenantsServer(grpcServer, systenantsServer.NewSysTenantsServer(ctx))
		apps.RegisterSysPackagesServer(grpcServer, syspackagesServer.NewSysPackagesServer(ctx))
		apps.RegisterSysApisServer(grpcServer, sysapisServer.NewSysApisServer(ctx))
		apps.RegisterSysLogsServer(grpcServer, syslogsServer.NewSysLogsServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
