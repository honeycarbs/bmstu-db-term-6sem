package apiserver

import "goodisgood/internal/storage"

type Config struct {
	BindAddr string `toml:"bind_addr"`
	LogLevel string `toml:"log_level"`

	storage *storage.Config `toml:"storage"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr: "8080",
		LogLevel: "debug",
		storage:  storage.NewConfig(),
	}
}
