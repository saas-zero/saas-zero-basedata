// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package main

import (
	"flag"
	"fmt"

	"github.com/saas-zero/saas-zero-basedata/api/internal/config"
	"github.com/saas-zero/saas-zero-basedata/api/internal/handler"
	"github.com/saas-zero/saas-zero-basedata/api/internal/middleware"
	"github.com/saas-zero/saas-zero-basedata/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/systemapis.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)

	server.Use(middleware.JwtAuth(c.JwtSecret, ctx.Redis))
	server.Use(middleware.CasbinAuth(ctx.Enforcer))
	server.Use(middleware.OperationLog(ctx.SysLogs))

	handler.RegisterHandlers(server, ctx)
	handler.RegisterInitRoutes(server, ctx)
	handler.RegisterLogRoutes(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
