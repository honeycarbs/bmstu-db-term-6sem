package apiserver

import "github.com/BurntSushi/toml"

type Config struct {
	BindAddr   string `toml:"bind_addr"`
	LogLevel   string `toml:"log_level"`
	SessionKey string `toml:"session_key"`

	DBURI      string `toml:"database_uri"`
	DBUsername string `toml:"database_username"`
	DBPassword string `toml:"database_password"`
}

type DBConfig struct {
	DBURI      string `toml:"database_uri"`
	DBUsername string `toml:"database_username"`
	DBPassword string `toml:"database_password"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr: "8080",
		LogLevel: "debug",
	}
}

func NewDBConfig(filename string) (*DBConfig, error) {
	var conf DBConfig
	_, err := toml.DecodeFile(filename, &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}
