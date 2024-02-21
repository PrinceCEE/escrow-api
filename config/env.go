package config

import "os"

type Env struct {
	PORT      string
	DSN       string
	REDIS_URL string
}

func newEnv() *Env {
	return &Env{
		PORT:      os.Getenv("PORT"),
		DSN:       os.Getenv("DSN"),
		REDIS_URL: os.Getenv("REDIS_URL"),
	}
}
