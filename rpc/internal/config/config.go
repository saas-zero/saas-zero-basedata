package config

import "github.com/zeromicro/go-zero/zrpc"

type PostgresConfig struct {
	DataSource string
}

type Config struct {
	zrpc.RpcServerConf
	Postgres PostgresConfig
}
