package test_config

import (
	"context"
	"os"
	"time"

	"github.com/Bupher-Co/bupher-api/config"
	"github.com/Bupher-Co/bupher-api/internal/repositories"
	"github.com/Bupher-Co/bupher-api/tests/utils/mocks/test_repositories"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type TestConfig struct {
	AuthRepository     repositories.IAuthRepository
	BusinessRepository repositories.IBusinessRepository
	EventRepository    repositories.IEventRepository
	OtpRepository      repositories.IOtpRepository
	TokenRepository    repositories.ITokenRepository
	UserRepository     repositories.IUserRepository
	DB                 *pgxpool.Pool
	RedisClient        *config.RedisClient
	Logger             *config.Logger
}

func NewTestConfig() *TestConfig {
	dbConfig, err := pgxpool.ParseConfig("postgres://postgres:password@localhost/bupher-test?sslmode=disable")
	if err != nil {
		panic(err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		panic(err)
	}

	timeout := 10 * time.Second
	return &TestConfig{
		Logger: config.NewLogger(
			zerolog.New(os.Stderr).Level(zerolog.DebugLevel).With().Timestamp().Logger(),
			zerolog.DebugLevel,
		),
		RedisClient:        &config.RedisClient{},
		DB:                 pool,
		AuthRepository:     test_repositories.NewAuthRepository(pool, timeout),
		BusinessRepository: test_repositories.NewBusinessRepository(pool, timeout),
		EventRepository:    test_repositories.NewEventRepository(pool, timeout),
		OtpRepository:      test_repositories.NewOtpRepository(pool, timeout),
		TokenRepository:    test_repositories.NewTokenRepository(pool, timeout),
		UserRepository:     test_repositories.NewUserRepository(pool, timeout),
	}
}

func (c *TestConfig) Getenv(key string) string {
	return os.Getenv(key)
}

func (c *TestConfig) GetAuthRepository() repositories.IAuthRepository {
	return c.AuthRepository
}

func (c *TestConfig) GetBusinessRepository() repositories.IBusinessRepository {
	return c.BusinessRepository
}

func (c *TestConfig) GetEventRepository() repositories.IEventRepository {
	return c.EventRepository
}

func (c *TestConfig) GetOtpRepository() repositories.IOtpRepository {
	return c.OtpRepository
}

func (c *TestConfig) GetTokenRepository() repositories.ITokenRepository {
	return c.TokenRepository
}

func (c *TestConfig) GetUserRepository() repositories.IUserRepository {
	return c.UserRepository
}

func (c *TestConfig) GetDB() *pgxpool.Pool {
	return c.DB
}

func (c *TestConfig) GetRedisClient() *config.RedisClient {
	return c.RedisClient
}

func (c *TestConfig) GetLogger() *config.Logger {
	return c.Logger
}
