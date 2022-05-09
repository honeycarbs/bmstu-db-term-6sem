package main

import (
	"api/internal/router"
	"api/pkg/logger"
)

func main() {
	logger.Init()
	logger := logger.GetLogger()
	logger.Println("Logger initialized.")

	// ctx := app_context.Singleton()
	logger.Println("Application context initialized.")

	defer router.Init()
	logger.Println("Application initialized. Have fun!")
}
