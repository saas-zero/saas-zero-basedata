package svc

import (
	"context"
	"fmt"

	"github.com/saas-zero/saas-zero-basedata/ent"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/config"

	"github.com/zeromicro/go-zero/core/logx"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

type ServiceContext struct {
	Config config.Config
	DB     *ent.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	log := logx.WithContext(context.Background())

	var err error = nil
	var client *ent.Client = nil

	// 添加调试信息
	log.Debugf("Database config: type=%s, host=%s, port=%d, user=%s, dbname=%s",
		c.Database.DBtype, c.Database.Host, c.Database.Port, c.Database.User, c.Database.DBName)

	switch c.Database.DBtype {
	case "postgres":
		dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s", c.Database.User, c.Database.Password, c.Database.Host,
			c.Database.Port, c.Database.DBName, c.Database.SSLMode)
		client, err = ent.Open(c.Database.DBtype, dsn)
	case "mysql":
		_ = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.Database.User,
			c.Database.Password, c.Database.Host, c.Database.Port, c.Database.DBName)
	default:
		log.Errorf("unsupported database type: %s", c.Database.DBtype)
	}

	if err != nil {
		log.Errorf("failed opening connection to postgres: %v", err)
	}
	defer client.Close()
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Errorf("failed creating schema resources: %v", err)
	}

	return &ServiceContext{
		Config: c,
		DB:     client,
	}
}
