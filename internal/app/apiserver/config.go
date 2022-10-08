package apiserver

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	BindAddr           string `toml:"bind_addr"`
	LogLevel           string `toml:"log_level"`
	DatabaseUrl        string `toml:"database_url"`
	DatabaseDriverName string `toml:"database_driver_name"`
	SessionKey         string `toml:"session_key"`
}

func NewConfigFromToml(path string) (*Config, error) {
	config := &Config{}
	if _, err := toml.DecodeFile(path, config); err != nil {
		return nil, err
	}
	return config, nil
}
