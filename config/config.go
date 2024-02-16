package config

import (
	"flag"
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	DbManager *DbManager
	Env       *Env
}

func NewConfig() *Config {
	var environment string
	flag.StringVar(&environment, "env", "development", "The environment of the app(development/production)")
	flag.Parse()

	if environment == "development" {
		if err := godotenv.Load(); err != nil {
			log.Panic(err)
		}
	}

	env := newEnv()
	return &Config{Env: env, DbManager: newDbManager(env)}
}
