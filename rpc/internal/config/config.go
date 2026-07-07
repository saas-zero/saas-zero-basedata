package config

import (
	"github.com/saas-zero/saas-zero-common/pkg/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type PostgresConfig struct {
	DataSource string
}

type Config struct {
	zrpc.RpcServerConf
	Postgres   PostgresConfig
	CacheRedis redis.Conf `json:"cacheRedis"`
}
