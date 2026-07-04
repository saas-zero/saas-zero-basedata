package main

import (
	"context"
	"fmt"
	"os"

	"entgo.io/ent/dialect"
	_ "github.com/lib/pq"
	"github.com/saas-zero/saas-zero-basedata/ent"
)

func main() {
	dsn := "postgresql://postgres:AM38xymTdFree4Fh@192.168.201.188:5432/saas_zero_kun?sslmode=disable"

	client, err := ent.Open(dialect.Postgres, dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed connecting: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	if err := client.Schema.Create(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "failed creating schema: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Migration completed successfully")
}
