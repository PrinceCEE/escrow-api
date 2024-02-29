package config

import (
	"context"
	"errors"
	"flag"
	"os"
	"time"

	"github.com/Bupher-Co/bupher-api/internal/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

type Logger struct {
	l     zerolog.Logger
	level zerolog.Level
}

type config struct {
	DB           *pgxpool.Pool
	Repositories Repositories
	Env          Env
	RedisClient  *redisClient
	Logger       Logger
}

type Repositories struct {
	UserRepository     repositories.UserRepository
	BusinessRepository repositories.BusinessRepository
	AuthRepository     repositories.AuthRepository
	EventRepository    repositories.EventRepository
	TokenRepository    repositories.TokenRepository
	OtpRepository      repositories.OtpRepository
}

type Env struct {
	PORT           string
	DSN            string
	REDIS_URL      string
	EMAIL_USERNAME string
	EMAIL_PASSWORD string
	ENVIRONMENT    string
	JWT_KEY        string
}

func (e *Env) IsDevelopment() bool {
	return e.ENVIRONMENT == "development"
}

func (e *Env) IsProduction() bool {
	return e.ENVIRONMENT == "production"
}

var Config = newConfig()

func newConfig() config {
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
	logger := Logger{
		l:     zerolog.New(os.Stderr).Level(level).With().Timestamp().Logger(),
		level: level,
	}

	env := Env{
		PORT:           os.Getenv("PORT"),
		DSN:            os.Getenv("DSN"),
		REDIS_URL:      os.Getenv("REDIS_URL"),
		EMAIL_USERNAME: os.Getenv("EMAIL_USERNAME"),
		EMAIL_PASSWORD: os.Getenv("EMAIL_PASSWORD"),
		ENVIRONMENT:    os.Getenv("ENVIRONMENT"),
		JWT_KEY:        os.Getenv("JWT_KEY"),
	}

	parsedConfig, err := pgxpool.ParseConfig(env.DSN)
	if err != nil {
		logger.Log(zerolog.PanicLevel, "error parsing dsn", nil, err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, parsedConfig)
	if err != nil {
		logger.Log(zerolog.PanicLevel, "error connecting the db", nil, err)
	}

	rclient, err := NewRedisClient(env)
	if err != nil {
		logger.Log(zerolog.PanicLevel, "error instantiating redis client", nil, err)
	}

	timeout := 10 * time.Second
	return config{
		DB:          pool,
		Env:         env,
		Logger:      logger,
		RedisClient: rclient,
		Repositories: Repositories{
			AuthRepository:     repositories.NewAuthRepository(pool, timeout),
			BusinessRepository: repositories.NewBusinessRepository(pool, timeout),
			EventRepository:    repositories.NewEventRepository(pool, timeout),
			TokenRepository:    repositories.NewTokenRepository(pool, timeout),
			UserRepository:     repositories.NewUserRepository(pool, timeout),
			OtpRepository:      repositories.NewOtpRepository(pool, timeout),
		},
	}
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
