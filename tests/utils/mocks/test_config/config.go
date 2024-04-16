package test_config

import (
	"context"
	"io"
	"time"

	"github.com/Bupher-Co/bupher-api/config"
	"github.com/Bupher-Co/bupher-api/internal/repositories"
	"github.com/Bupher-Co/bupher-api/pkg/apis"
	"github.com/Bupher-Co/bupher-api/pkg/apis/paystack"
	"github.com/Bupher-Co/bupher-api/pkg/push"
	"github.com/Bupher-Co/bupher-api/tests/utils/mocks/test_repositories"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
)

type TestConfig struct {
	AuthRepository                repositories.IAuthRepository
	BusinessRepository            repositories.IBusinessRepository
	EventRepository               repositories.IEventRepository
	OtpRepository                 repositories.IOtpRepository
	TokenRepository               repositories.ITokenRepository
	UserRepository                repositories.IUserRepository
	WalletRepository              repositories.IWalletRepository
	WalletHistoryRepository       repositories.IWalletHistoryRepository
	BankAccountRepository         repositories.IBankAccountRepository
	TransactionRepository         repositories.ITransactionRepository
	TransactionTimelineRepository repositories.ITransactionTimelineRepository
	DB                            *pgxpool.Pool
	RedisClient                   *config.RedisClient
	Logger                        *config.Logger
	Push                          push.IPush
	Apis                          apis.IAPIs
	mock.Mock
}

type TestPush struct{ mock.Mock }

func (p *TestPush) SendEmail(data *push.Email) error {
	return nil
}

func (p *TestPush) SendSMS(data *push.Sms) {}

type TestAPIs struct{ mock.Mock }

type TestPaystackAPI struct{ mock.Mock }

func (a *TestAPIs) GetPaystack() paystack.IPaystack {
	return &TestPaystackAPI{}
}

func (p *TestPaystackAPI) InitiateTransaction(data paystack.InitiateTransactionDto) (*paystack.InitiateTransactionResponse, error) {
	return &paystack.InitiateTransactionResponse{}, nil
}

func NewTestConfig() *TestConfig {
	dbConfig, err := pgxpool.ParseConfig("postgres://postgres:password@localhost/bupher_test?sslmode=disable")
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
			zerolog.New(io.Discard).Level(zerolog.DebugLevel).With().Timestamp().Logger(),
			zerolog.DebugLevel,
		),
		RedisClient:                   &config.RedisClient{},
		DB:                            pool,
		AuthRepository:                test_repositories.NewAuthRepository(pool, timeout),
		BusinessRepository:            test_repositories.NewBusinessRepository(pool, timeout),
		EventRepository:               test_repositories.NewEventRepository(pool, timeout),
		OtpRepository:                 test_repositories.NewOtpRepository(pool, timeout),
		TokenRepository:               test_repositories.NewTokenRepository(pool, timeout),
		UserRepository:                test_repositories.NewUserRepository(pool, timeout),
		WalletRepository:              test_repositories.NewWalletRepository(pool, timeout),
		WalletHistoryRepository:       test_repositories.NewWalletHistoryRepository(pool, timeout),
		BankAccountRepository:         test_repositories.NewBankAccountRepository(pool, timeout),
		TransactionRepository:         test_repositories.NewTransactionRepository(pool, timeout),
		TransactionTimelineRepository: test_repositories.NewTransactionTimelineRepository(pool, timeout),
		Push:                          &TestPush{},
	}
}

func (c *TestConfig) Getenv(key string) string {
	switch key {
	case "ENVIRONMENT":
		return "test"
	case "JWT_KEY":
		return "somerandomjwtkey"
	default:
		return ""
	}
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

func (c *TestConfig) GetWalletRepository() repositories.IWalletRepository {
	return c.WalletRepository
}

func (c *TestConfig) GetWalletHistoryRepository() repositories.IWalletHistoryRepository {
	return c.WalletHistoryRepository
}

func (c *TestConfig) GetBankAccountRepository() repositories.IBankAccountRepository {
	return c.BankAccountRepository
}

func (c *TestConfig) GetTransactionRepository() repositories.ITransactionRepository {
	return c.TransactionRepository
}

func (c *TestConfig) GetTransactionTimelineRepository() repositories.ITransactionTimelineRepository {
	return c.TransactionTimelineRepository
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

func (c *TestConfig) GetPush() push.IPush {
	return &push.Push{}
}

func (c *TestConfig) GetAPIs() apis.IAPIs {
	return &TestAPIs{}
}
