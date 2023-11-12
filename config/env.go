package config

import "os"

type Env struct {
	PORT string
}

func newEnv() *Env {
	return &Env{
		PORT: os.Getenv("PORT"),
	}
}
