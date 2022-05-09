package config

import (
	"api/pkg/logger"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	isDebug *bool `yaml:"is_debug"`
	Listen  struct {
		Type   string `yaml:"type" env-default:"port"`
		BindIP string `yaml:"bind_ip" env-default:"localhost"`
		Port   string `yaml:"port" env-default:"8080"`
	}
}

var instance *Config

func initConfig() {
	logger := logger.GetLogger()
	logger.Info("Read application config")
	instance = &Config{}
	if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
		ret, _ := cleanenv.GetDescription(instance, nil)
		logger.Info(ret)
		logger.Fatal(err)
	}
}

func GetConfig() *Config {
	if instance == nil {
		initConfig()
	}
	return instance
}
