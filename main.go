package main

import (
	"context"
	"fmt"
	"river-args-reproduce/worker"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
)

type api struct {
	client *river.Client[pgx.Tx]
	pool   *pgxpool.Pool
}

func main() {
	ctx := context.TODO()

	config, err := pgxpool.ParseConfig(
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			"postgres", "postgres", "localhost", "5432", "postgres"),
	)
	if err != nil {
		panic(err)
	}

	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	config.MaxConnLifetime = 1 * time.Hour
	config.MaxConnIdleTime = 30 * time.Second

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		panic(err)
	}

	workers := river.NewWorkers()
	river.AddWorker(workers, worker.NewExampleWorker())

	client, err := river.NewClient[pgx.Tx](riverpgxv5.New(pool), &river.Config{
		Workers: workers,
	})
	if err != nil {
		panic(err)
	}

	a := &api{client: client, pool: pool}

	err = a.insertJob()
	if err != nil {
		panic(err)
	}
}

func (a *api) insertJob() error {
	ctx := context.TODO()

	_, err := a.client.Insert(ctx, worker.ExampleJobArgs{}, nil)
	if err != nil {
		return err
	}

	return nil
}
