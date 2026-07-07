// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type CasbinPostgresConfig struct {
	DataSource string `json:"dataSource"`
}

type Config struct {
	rest.RestConf
	JwtSecret      string               `json:"jwtSecret"`
	Redis          redis.RedisConf      `json:"redis"`
	CasbinPostgres CasbinPostgresConfig `json:"casbinPostgres"`
	Basedata       zrpc.RpcClientConf
}
