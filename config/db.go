package config

import (
	"context"
	"time"

	"github.com/Bupher-Co/bupher-api/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositories struct {
	UserRepository     repositories.UserRepository
	BusinessRepository repositories.BusinessRepository
	AuthRepository     repositories.AuthRepository
	EventRepository    repositories.EventRepository
	TokenRepository    repositories.TokenRepository
}

type DbManager struct {
	DB           *pgxpool.Pool
	Repositories Repositories
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

	timeout := 10 * time.Second
	return &DbManager{
		DB: pool,
		Repositories: Repositories{
			AuthRepository:     repositories.NewAuthRepository(pool, timeout),
			BusinessRepository: repositories.NewBusinessRepository(pool, timeout),
			EventRepository:    repositories.NewEventRepository(pool, timeout),
			TokenRepository:    repositories.NewTokenRepository(pool, timeout),
			UserRepository:     repositories.NewUserRepository(pool, timeout),
		},
	}, nil
}
