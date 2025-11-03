package main

import (
	"flag"
	"fmt"

	"github.com/saas-zero/saas-zero-basedata/rpc/apps/system-service"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/config"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/server"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "D:\\GolandProjects\\saas-zero\\apps\\saas-zero-basedata\\rpc\\etc\\systemservice.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		system_service.RegisterSysUsersServer(grpcServer, server.NewSysUsersServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
