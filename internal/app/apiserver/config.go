package apiserver

import "awesomeProject/internal/app/store"

type Config struct {
	BindAddr string `toml:"bind_addr"`
	LogLevel string `toml:"log_level"`
	Store    *store.Config
}

func NewDefaultConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogLevel: "Debug",
		Store:    store.NewConfig(),
	}
}
