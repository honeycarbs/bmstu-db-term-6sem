package apiserver

type Config struct {
	BindAddr   string `toml:"bind_addr"`
	LogLevel   string `toml:"log_level"`
	DBURI      string `toml:"database_uri"`
	DBUsername string `toml:"database_username"`
	DBPassword string `toml:"database_password"`
	SessionKey string `toml:"session_key"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr: "8080",
		LogLevel: "debug",
	}
}
