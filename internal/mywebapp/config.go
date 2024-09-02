package mywebapp

import "github.com/VeeRomanoff/mywebapp/storage"

type Config struct {
	Port          string          `toml:"server_port"`
	Logger        string          `toml:"logger_level"`
	StorageConfig *storage.Config `toml:"storage"`
}

func NewConfig() *Config {
	return &Config{
		// default values
		Port:          ":8080",
		Logger:        "debug",
		StorageConfig: storage.New(),
	}
}
