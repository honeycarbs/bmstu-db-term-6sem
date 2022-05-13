package storage

type Config struct {
	DBuri      string `toml:"database_uri"`
	DBUsername string `toml:"database_username"`
	DBPassword string `toml:"database_password"`
}

func NewConfig() *Config {
	return &Config{
		DBuri:      "bolt://localhost",
		DBUsername: "neo4j",
		DBPassword: "test",
	}
}
