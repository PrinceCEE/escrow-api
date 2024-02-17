package config

import (
	"errors"
	"flag"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

type Logger struct {
	l     zerolog.Logger
	level zerolog.Level
}

type Config struct {
	DbManager *DbManager
	Env       *Env
	Logger    *Logger
}

var Cfg = newConfig()

func newConfig() *Config {
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
	logger := &Logger{
		l:     zerolog.New(os.Stderr).Level(level).With().Timestamp().Logger(),
		level: level,
	}

	env := newEnv()
	manager, err := newDbManager(env)
	if err != nil {
		logger.Log(zerolog.PanicLevel, "error connecting to the db", nil, err)
	}

	return &Config{
		Env:       env,
		DbManager: manager,
		Logger:    logger,
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
