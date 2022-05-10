package main

import (
	"api/internal/config"
	"api/internal/handlers/auth"
	"api/internal/handlers/landing"
	"api/internal/handlers/poll"
	"api/pkg/cache/freecache"
	"api/pkg/handler/metric"
	"api/pkg/logger"
	"api/pkg/shutdown"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	logger.Init()
	logger := logger.GetLogger()
	logger.Println("Logger initialized.")

	logger.Println("Config initialized.")
	cfg := config.GetConfig()

	router := httprouter.New()
	logger.Println("Router initialized.")

	logger.Println("RegisterCache initialized.")
	refreshTokenCache := freecache.NewCacheRepository(104857600) // 100 MB

	authHandler := auth.Handler{RTCache: refreshTokenCache, Logger: logger}
	authHandler.Register(router)

	metricHandler := metric.Handler{Logger: logger}
	metricHandler.Register(router)

	pollHandler := poll.Handler{Logger: logger}
	pollHandler.Register(router)

	landingHandler := landing.Handler{Logger: logger}
	landingHandler.Register(router)

	logger.Println("Starting application....")
	start(router, logger, cfg)
}

func start(router *httprouter.Router, logger logger.Logger, cfg *config.Config) {
	var server *http.Server
	var listener net.Listener

	if cfg.Listen.Type == "sock" {
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}
		socketPath := path.Join(appDir, "app.sock")
		logger.Info("Socket path: %s", socketPath)

		listener, err = net.Listen("unix", socketPath)
		if err != nil {
			logger.Fatal(err)
		}
	} else {
		logger.Infof("Bind application host: %s, port: %s",
			cfg.Listen.BindIP, cfg.Listen.Type)
		var err error
		listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s",
			cfg.Listen.BindIP, cfg.Listen.Port))
		if err != nil {
			logger.Fatal(err)
		}
	}

	server = &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go shutdown.Graceful([]os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGTSTP, syscall.SIGHUP,
		syscall.SIGTERM, os.Interrupt, os.Kill}, server)

	if err := server.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			logger.Warn("Server shutdown")
		default:
			logger.Fatal(err)
		}
	}
	logger.Println("Application initialized and starterd.")

}
