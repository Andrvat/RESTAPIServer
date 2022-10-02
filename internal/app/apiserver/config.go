package apiserver

type Config struct {
	BindAddr string `toml:"bind_addr"`
	LogLevel string `toml:"log_level"`
}

func NewDefaultConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogLevel: "Debug",
	}
}
