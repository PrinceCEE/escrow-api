package config

import "os"

type Env struct {
	PORT string
	DSN  string
}

func newEnv() *Env {
	return &Env{
		PORT: os.Getenv("PORT"),
		DSN:  os.Getenv("DSN"),
	}
}
