package app_context

import (
	"api/internal/config"
	"api/pkg/logger"
	"sync"
)

type AppContext struct {
	Config *config.Config
}

var instance *AppContext
var once sync.Once

func Singleton() *AppContext {
	logger := logger.GetLogger()

	logger.Println("Initializing application context")
	once.Do(func() {
		instance = &AppContext{
			Config: config.GetConfig(),
		}
	})

	return instance
}
