package config

import (
	"context"
	"time"

	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func configureDB(dsn string) (*pgxpool.Pool, error) {
	parsedConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	setupHooks(parsedConfig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, parsedConfig)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func setupHooks(parseConfig *pgxpool.Config) {
	parseConfig.AfterConnect = func(ctx context.Context, c *pgx.Conn) error {
		pgxuuid.Register(c.TypeMap())
		return nil
	}
}
