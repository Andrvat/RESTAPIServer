package apiserver

type Config struct {
	BindAddr           string `toml:"bind_addr"`
	LogLevel           string `toml:"log_level"`
	DatabaseUrl        string `toml:"database_url"`
	DatabaseDriverName string `toml:"database_driver_name"`
}

func NewDefaultConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogLevel: "Debug",
	}
}
