package store

type Config struct {
	DatabaseUrl        string `toml:"database-url"`
	DatabaseDriverName string `toml:"database-driver-name"`
}

func NewConfig() *Config {
	return &Config{}
}
