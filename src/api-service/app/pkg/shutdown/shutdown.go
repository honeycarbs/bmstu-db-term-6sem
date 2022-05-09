package shutdown

import (
	"api/pkg/logger"
	"io"
	"os"
	"os/signal"
)

func Graceful(signals []os.Signal, closeItems ...io.Closer) {
	logger := logger.GetLogger()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, signals...)
	sig := <-sigc
	logger.Infof("Caught signal %s. Shutting down.", sig)

	for _, closer := range closeItems {
		if err := closer.Close(); err != nil {
			logger.Errorf("failed to close %v: %v", closer, err)
		}
	}
	// go func(c chan os.Signal) {
	// 	sig := <-c
	// 	logger.Info("Caught signal %s. Shutting down.", sig)
	// 	err := unixListener.Close()
	// 	if err != nil {
	// 		logger.Error("Can't close application unix socket")
	// 	}
	// 	os.Exit(0)
	// }(sigc)
}
