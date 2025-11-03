// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"system-service/api/internal/config"
	system_service "system-service/rpc/apps/system-service"
)

type ServiceContext struct {
	Config        config.Config
	SystemService system_service.SysUsersClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		SystemService: system_service.SysUsersClient.NewSysUsers(c.SystemRpc),
	}
}
