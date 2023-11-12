package config

import (
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	DbManager *DbManager
	Env       *Env
}

func NewConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Panic(err)
	}

	env := newEnv()
	return &Config{Env: env, DbManager: newDbManager(env)}
}
