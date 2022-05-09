package router

import (
	"api/internal/app_context"
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

	h "api/pkg/metric/handler"

	"github.com/julienschmidt/httprouter"
)

func Init() {
	logger := logger.GetLogger()
	logger.Println("Initialize application router")

	router := httprouter.New()
	router.GET(h.HEARTBEET_URL, h.Heartbeet)

	var server *http.Server
	var listener net.Listener

	ctx := app_context.Singleton()
	cfg := ctx.Config

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

	go shutdown.Graceful([]os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, syscall.SIGTERM, os.Interrupt})

	if err := server.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			logger.Warn("Server shutdown")
		default:
			logger.Fatal(err)
		}
	}
}
