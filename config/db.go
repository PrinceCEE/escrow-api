package config

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DbManager struct {
	DB *pgxpool.Pool
}

func newDbManager(env *Env) (*DbManager, error) {
	config, err := pgxpool.ParseConfig(env.DSN)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return &DbManager{DB: pool}, nil
}

func (m *DbManager) Save() interface{} {
	return "Not implemented"
}

func (m *DbManager) Update() interface{} {
	return "Not Implemented"
}

func (m *DbManager) FindOne() interface{} {
	return "Not implemented"
}

func (m *DbManager) Find() interface{} {
	return "Not implemented"
}

func (m *DbManager) SoftDelete() interface{} {
	return "Not implemented"
}

func (m *DbManager) HardDelete() interface{} {
	return "Not implemented"
}
