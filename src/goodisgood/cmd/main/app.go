package main

import (
	"flag"
	"goodisgood/internal/app/apiserver"

	// "github.com/sirupsen/logrus"
	"log"

	"github.com/BurntSushi/toml"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "config/config.toml", "path to config file")
}

func main() {
	flag.Parse()

	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	// s := apiserver.Start(config)
	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
