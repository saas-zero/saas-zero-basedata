package svc

import (
	"context"
	"database/sql"
	"log"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"
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
	client = client.Debug()
	// WithDropIndex: 删除 schema 中已不存在的索引（如旧的字段级唯一索引），
	// 确保条件唯一索引迁移后旧索引不会残留导致唯一约束仍生效。
	if err := client.Schema.Create(context.Background(), schema.WithDropIndex(true)); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	if err := SeedSystemDicts(context.Background(), client); err != nil {
		log.Fatalf("failed seeding system dictionaries: %v", err)
	}
	casbinDb, err := sql.Open("postgres", c.Postgres.DataSource)
	if err != nil {
		log.Fatalf("failed opening casbin db: %v (fail-closed)", err)
	}
	enf, err := commcasbin.NewEnforcer(casbinDb, "casbin_rule")
	if err != nil {
		log.Fatalf("failed initializing casbin: %v (fail-closed)", err)
	}
	rds, err := redis.NewClient(c.CacheRedis)
	if err != nil {
		log.Fatalf("failed initializing redis: %v (fail-closed)", err)
	}
	return &ServiceContext{
		Config:   c,
		DB:       client,
		Redis:    rds,
		Enforcer: enf,
	}
}
