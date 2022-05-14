package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"goodisgood/internal/app/model"
	"goodisgood/internal/app/storage"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

const (
	sessionName                  = "goodisgood"
	contextKeyAccount contextKey = iota
	contextKeyRequestID
)

var (
	errInvalidEmailOrPassword = errors.New("invalid e-mail or password.")
	errAccountAlreadyExists   = errors.New("user already exists.")
	errNotAuthentificated     = errors.New("Not authentificated")
)

type contextKey int8

type server struct {
	router         *mux.Router
	logger         *logrus.Logger
	storage        storage.Storage
	sessionStorage sessions.Store
}

func newServer(storage storage.Storage, sessionStorage sessions.Store) *server {
	s := &server{
		router:         mux.NewRouter(),
		logger:         logrus.New(),
		storage:        storage,
		sessionStorage: sessionStorage,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	s.router.HandleFunc("/accounts", s.handleAccountsCreate()).Methods("POST")
	s.router.HandleFunc("/sessions", s.handleSessionsCreate()).Methods("POST")

	private := s.router.PathPrefix("/private").Subrouter()
	private.Use(s.authenticateAccount)
	private.HandleFunc("/whoami", s.handleWhoami()).Methods("GET")
}

func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), contextKeyRequestID, id)))
	})
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(contextKeyRequestID),
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		logger.Infof("completed with %d %s in: %v",
			rw.statusCode, http.StatusText(rw.statusCode),
			time.Now().Sub(start))
	})
}

func (s *server) authenticateAccount(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStorage.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		uuid, ok := session.Values["account_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthentificated)
		}

		a, err := s.storage.Account().Find(uuid.(string))
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthentificated)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), contextKeyAccount, a)))
	})
}

func (s *server) handleWhoami() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, r.Context().Value(contextKeyAccount).(*model.Account))
	}
}

func (s *server) handleAccountsCreate() http.HandlerFunc {
	type request struct {
		UUID     string `json:"account_uuid"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		User     struct {
			UUUID  string `json:"uuid"`
			Age    int    `json:"age"`
			Race   string `json:"race"`
			Gender string `json:"gender"`
		} `json:"user"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.logger.Error(w, r, http.StatusBadRequest, err)
			return
		}
		s.logger.Info(req)

		a := &model.Account{
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password,
		}

		if _, err := s.storage.Account().FindByEmail(a.Email); err == nil {
			s.error(w, r, http.StatusUnauthorized, errAccountAlreadyExists)
			return
		}

		if err := s.storage.Account().Create(a); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		u := &model.User{
			UUID:   req.User.UUUID,
			Age:    req.User.Age,
			Race:   req.User.Race,
			Gender: req.User.Gender,
		}

		if err := s.storage.User().Create(a.UUID, u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		a.Sanitize()
		s.respond(w, r, http.StatusCreated, a)
	}
}

func (s *server) handleSessionsCreate() http.HandlerFunc {
	type request struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.logger.Error(w, r, http.StatusBadRequest, err)
			return
		}
		a, err := s.storage.Account().FindByEmail(req.Email)
		if err != nil || !a.ComparePassword(req.Password) {
			if err != nil {
				s.logger.Infof(err.Error())
			}
			s.error(w, r, http.StatusUnauthorized, errInvalidEmailOrPassword)
			return
		}

		session, err := s.sessionStorage.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Values["account_id"] = a.UUID
		if err := s.sessionStorage.Save(r, w, session); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
