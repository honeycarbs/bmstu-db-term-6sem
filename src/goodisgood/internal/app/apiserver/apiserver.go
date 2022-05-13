package apiserver

import (
	"goodisgood/internal/storage"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type APIServer struct {
	config  *Config
	logger  *logrus.Logger
	router  *mux.Router
	storage *storage.Storage
}

func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

func (s *APIServer) Start() error {
	if err := s.getLogger(); err != nil {
		return err
	}

	s.getRouter()

	if err := s.getStorage(); err != nil {
		return err
	}

	s.logger.Info("Starting API server...")

	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *APIServer) getLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}

	s.logger.SetLevel(level)

	return nil
}

func (s *APIServer) getRouter() {
	s.router.HandleFunc("/hello", s.handleHello())
}

func (s *APIServer) getStorage() error {
	st := storage.New(s.config.storage)
	if err := st.Open(); err != nil {
		return err
	}
	s.storage = st

	return nil
}

func (s *APIServer) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hello")
	}
}
