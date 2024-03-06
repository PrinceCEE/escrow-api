package config

import (
	"errors"
	"flag"
	"os"
	"time"

	"github.com/Bupher-Co/bupher-api/internal/repositories"
	"github.com/Bupher-Co/bupher-api/pkg/push"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

type Logger struct {
	l     zerolog.Logger
	level zerolog.Level
}

func NewLogger(l zerolog.Logger, level zerolog.Level) *Logger {
	return &Logger{l, level}
}

func getLoggerLevel(loglevel string) zerolog.Level {
	switch loglevel {
	case zerolog.LevelTraceValue:
		return zerolog.TraceLevel

	case zerolog.LevelDebugValue:
		return zerolog.DebugLevel

	case zerolog.LevelInfoValue:
		return zerolog.InfoLevel

	case zerolog.LevelWarnValue:
		return zerolog.WarnLevel

	case zerolog.LevelErrorValue:
		return zerolog.ErrorLevel

	case zerolog.LevelFatalValue:
		return zerolog.FatalLevel

	case zerolog.LevelPanicValue:
		return zerolog.PanicLevel
	default:
		panic(errors.New("invalid loglevel value"))
	}
}

func (logger *Logger) Log(level zerolog.Level, msg string, data map[string]any, err error) {
	var lEvent *zerolog.Event

	switch level {
	case zerolog.TraceLevel:
		lEvent = logger.l.Trace()

	case zerolog.DebugLevel:
		lEvent = logger.l.Debug()

	case zerolog.InfoLevel:
		lEvent = logger.l.Info()

	case zerolog.WarnLevel:
		lEvent = logger.l.Warn()

	case zerolog.ErrorLevel:
		lEvent = logger.l.Error()

	case zerolog.FatalLevel:
		lEvent = logger.l.Fatal()

	case zerolog.PanicLevel:
		lEvent = logger.l.Panic()
	}

	for k, v := range data {
		lEvent = lEvent.Any(k, v)
	}

	if err != nil {
		lEvent = lEvent.Err(err)
	}

	lEvent.Msg(msg)
}

type IConfig interface {
	Getenv(key string) string
	GetAuthRepository() repositories.IAuthRepository
	GetBusinessRepository() repositories.IBusinessRepository
	GetEventRepository() repositories.IEventRepository
	GetOtpRepository() repositories.IOtpRepository
	GetTokenRepository() repositories.ITokenRepository
	GetUserRepository() repositories.IUserRepository
	GetDB() *pgxpool.Pool
	GetRedisClient() *RedisClient
	GetLogger() *Logger
	GetPush() push.IPush
}

type Config struct {
	AuthRepository     repositories.IAuthRepository
	BusinessRepository repositories.IBusinessRepository
	EventRepository    repositories.IEventRepository
	OtpRepository      repositories.IOtpRepository
	TokenRepository    repositories.ITokenRepository
	UserRepository     repositories.IUserRepository
	DB                 *pgxpool.Pool
	RedisClient        *RedisClient
	Logger             *Logger
	Push               push.IPush
}

func NewConfig() *Config {
	var environment, loglevel string

	flag.StringVar(&environment, "env", "development", "The environment of the app(development/production)")
	flag.StringVar(&loglevel, "loglevel", "debug", "The logger log level")
	flag.Parse()

	if environment == "development" {
		if err := godotenv.Load(); err != nil {
			panic(err)
		}
	}

	level := getLoggerLevel(loglevel)
	logger := NewLogger(
		zerolog.New(os.Stderr).Level(level).With().Timestamp().Logger(),
		level,
	)

	dbpool, err := configureDB(os.Getenv("DSN"))
	if err != nil {
		logger.Log(zerolog.PanicLevel, "error connecting to the db", nil, err)
	}

	rclient, err := NewRedisClient(os.Getenv("REDIS_URL"))
	if err != nil {
		logger.Log(zerolog.PanicLevel, "error instantiating redis client", nil, err)
	}

	timeout := 10 * time.Second
	return &Config{
		DB:                 dbpool,
		Logger:             logger,
		RedisClient:        rclient,
		AuthRepository:     repositories.NewAuthRepository(dbpool, timeout),
		BusinessRepository: repositories.NewBusinessRepository(dbpool, timeout),
		EventRepository:    repositories.NewEventRepository(dbpool, timeout),
		TokenRepository:    repositories.NewTokenRepository(dbpool, timeout),
		UserRepository:     repositories.NewUserRepository(dbpool, timeout),
		OtpRepository:      repositories.NewOtpRepository(dbpool, timeout),
		Push:               &push.Push{},
	}
}

func (c *Config) Getenv(key string) string {
	return os.Getenv(key)
}

func (c *Config) GetAuthRepository() repositories.IAuthRepository {
	return c.AuthRepository
}

func (c *Config) GetBusinessRepository() repositories.IBusinessRepository {
	return c.BusinessRepository
}

func (c *Config) GetEventRepository() repositories.IEventRepository {
	return c.EventRepository
}

func (c *Config) GetOtpRepository() repositories.IOtpRepository {
	return c.OtpRepository
}

func (c *Config) GetTokenRepository() repositories.ITokenRepository {
	return c.TokenRepository
}

func (c *Config) GetUserRepository() repositories.IUserRepository {
	return c.UserRepository
}

func (c *Config) GetDB() *pgxpool.Pool {
	return c.DB
}

func (c *Config) GetRedisClient() *RedisClient {
	return c.RedisClient
}

func (c *Config) GetLogger() *Logger {
	return c.Logger
}

func (c *Config) GetPush() push.IPush {
	return &push.Push{}
}
