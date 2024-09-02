package storage

type Config struct {
	DatabaseURI string `toml:"database_uri"`
}

func New() *Config {
	return &Config{}
}
