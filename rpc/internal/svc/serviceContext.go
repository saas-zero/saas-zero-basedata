package svc

import (
	"context"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/config"
	"log"

	"entgo.io/ent/dialect"
	_ "github.com/lib/pq"
	"github.com/saas-zero/saas-zero-basedata/ent"
	_ "github.com/saas-zero/saas-zero-basedata/ent/runtime"
)

type ServiceContext struct {
	Config config.Config
	DB     *ent.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	client, err := ent.Open(dialect.Postgres, c.Postgres.DataSource)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	return &ServiceContext{
		Config: c,
		DB:     client,
	}
}
