package svc

import (
	"context"
	"database/sql"
	"log"

	"entgo.io/ent/dialect"
	casbinapi "github.com/casbin/casbin/v2"
	_ "github.com/lib/pq"
	"github.com/saas-zero/saas-zero-basedata/ent"
	_ "github.com/saas-zero/saas-zero-basedata/ent/runtime"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/config"
	commcasbin "github.com/saas-zero/saas-zero-common/pkg/casbin"
	"github.com/saas-zero/saas-zero-common/pkg/redis"
)

type ServiceContext struct {
	Config   config.Config
	DB       *ent.Client
	Redis    *redis.Client
	Enforcer *casbinapi.SyncedEnforcer
}

func NewServiceContext(c config.Config) *ServiceContext {
	client, err := ent.Open(dialect.Postgres, c.Postgres.DataSource)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	casbinDb, err := sql.Open("postgres", c.Postgres.DataSource)
	if err != nil {
		log.Fatalf("failed opening casbin db: %v", err)
	}
	enf, err := commcasbin.NewEnforcer(casbinDb, "casbin_rule")
	if err != nil {
		log.Fatalf("failed initializing casbin: %v", err)
	}
	rds, err := redis.NewClient(c.CacheRedis)
	if err != nil {
		log.Fatalf("failed initializing redis: %v", err)
	}
	return &ServiceContext{
		Config:   c,
		DB:       client,
		Redis:    rds,
		Enforcer: enf,
	}
}
