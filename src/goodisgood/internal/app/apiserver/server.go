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
	errNotEnoughRights        = errors.New("Not enough rights to visit this page")
)

type contextKey int8

type server struct {
	router         *mux.Router
	logger         *logrus.Logger
	ustorage       storage.Storage
	astorage       storage.Storage
	mstorage       storage.Storage
	sessionStorage sessions.Store
}

func newServer(u, m, a storage.Storage, sessionStorage sessions.Store) *server {
	s := &server{
		router:         mux.NewRouter(),
		logger:         logrus.New(),
		ustorage:       u,
		astorage:       a,
		mstorage:       m,
		sessionStorage: sessionStorage,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	// s.router.Use(s.setRequestID)
	// s.router.Use(s.logRequest)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))

	s.router.HandleFunc("/accounts", s.handleAccountsCreate()).Methods("POST")
	s.router.HandleFunc("/sessions", s.handleSessionsCreate()).Methods("POST")
	s.router.HandleFunc("/locations", s.handleLocationGet()).Methods("GET")
	s.router.HandleFunc("/education", s.handleEducationPlaceGet()).Methods("GET")

	authorized := s.router.PathPrefix("/authorized").Subrouter()
	authorized.Use(s.authenticateAccount)
	authorized.HandleFunc("/whoami", s.handleWhoami()).Methods("GET")
	authorized.HandleFunc("/poll", s.handleSubmitPoll()).Methods("POST")
	// authorized.HandleFunc("/poll", s.handleResultGet()).Methods("GET")
	authorized.HandleFunc("/poll", s.handleWordsGet()).Methods("GET")

	moderate := s.router.PathPrefix("/moderate").Subrouter()
	moderate.Use(s.checkIfModerator)
	moderate.HandleFunc("/accounts", s.handleAllAccountsGet()).Methods("GET")

	administrate := s.router.PathPrefix("/administrate").Subrouter()
	administrate.Use(s.checkIfAdmin)
	administrate.HandleFunc("/stats", s.handleStatsGet()).Methods("GET")

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

		a, err := s.ustorage.Account().Find(uuid.(string))
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthentificated)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), contextKeyAccount, a)))
	})
}

func (s *server) checkIfModerator(next http.Handler) http.Handler {
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

		a, err := s.ustorage.Account().Find(uuid.(string))
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthentificated)
			return
		}
		s.logger.Info(a.Role)

		if a.Role != "moderator" {
			s.error(w, r, http.StatusUnauthorized, errNotEnoughRights)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), contextKeyAccount, a)))
	})
}

func (s *server) checkIfAdmin(next http.Handler) http.Handler {
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

		a, err := s.ustorage.Account().Find(uuid.(string))
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthentificated)
			return
		}
		s.logger.Info(a.Role)

		if a.Role != "admin" {
			s.error(w, r, http.StatusUnauthorized, errNotEnoughRights)
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

func (s *server) handleAddLocation() http.HandlerFunc {
	type request struct {
		Locations []struct {
			Name     string `json:"name"`
			Region   string `json:"region"`
			District string `json:"district"`
		} `json:"locations"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.logger.Error(w, r, http.StatusBadRequest, err)
			return
		}

		acc := r.Context().Value(contextKeyAccount).(*model.Account)

		for _, i := range req.Locations {
			l := &model.Location{
				Name:     i.Name,
				Region:   i.Region,
				District: i.District,
			}
			if err := s.ustorage.Location().Create(l); err != nil {
				s.logger.Error(w, r, http.StatusBadRequest, err)
				return
			}

			if err := s.ustorage.Location().Assign(acc.UUID, l); err != nil {
				s.logger.Error(w, r, http.StatusBadRequest, err)
				return
			}
		}
	}
}

func (s *server) handleAccountsCreate() http.HandlerFunc {
	type request struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		User     struct {
			Age    int    `json:"age"`
			Race   string `json:"race"`
			Gender string `json:"gender"`
		} `json:"user"`
		Locations []struct {
			Name     string `json:"name"`
			Region   string `json:"region"`
			District string `json:"district"`
		} `json:"locations"`
		EducationPlace []struct {
			Field            string `json:"field"`
			Level            string `json:"level"`
			Name             string `json:"name"`
			EducationProgram struct {
				Field string `json:"field"`
				Level string `json:"level"`
			} `json:"education_program"`
		} `json:"education_place"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.logger.Error(w, r, http.StatusBadRequest, err)
			return
		}

		a := &model.Account{
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password,
		}

		if _, err := s.ustorage.Account().FindByEmail(a.Email); err == nil {
			s.error(w, r, http.StatusUnauthorized, errAccountAlreadyExists)
			return
		}

		if err := s.ustorage.Account().Create(a); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		u := &model.User{
			UUID:   a.UUID,
			Age:    req.User.Age,
			Race:   req.User.Race,
			Gender: req.User.Gender,
		}

		if err := s.ustorage.User().Create(a.UUID, u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		for _, i := range req.Locations {
			l := &model.Location{
				Name:     i.Name,
				Region:   i.Region,
				District: i.District,
			}
			if err := s.ustorage.Location().Create(l); err != nil {
				s.logger.Error(w, r, http.StatusBadRequest, err)
				return
			}

			if err := s.ustorage.Location().Assign(a.UUID, l); err != nil {
				s.logger.Error(w, r, http.StatusBadRequest, err)
				return
			}
		}

		for _, i := range req.EducationPlace {
			e := &model.EducationPlace{
				Name: i.Name,
			}
			p := &model.EducationProgram{
				Field: i.EducationProgram.Field,
				Level: i.EducationProgram.Level,
			}
			if err := s.ustorage.EducationPlace().Assign(a.UUID, e, p); err != nil {
				s.logger.Error(w, r, http.StatusBadRequest, err)
				return
			}
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
		a, err := s.ustorage.Account().FindByEmail(req.Email)
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

func (s *server) handleLocationGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		larr, err := s.ustorage.Location().Get()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusCreated, larr)
	}
}

func (s *server) handleWordsGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := s.ustorage.Poll().GetWordsList()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusCreated, res)
	}
}

func (s *server) handleResultGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		acc := r.Context().Value(contextKeyAccount).(*model.Account)

		res, err := s.ustorage.Poll().GetUserResult(acc.UUID)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusCreated, res)
	}
}

func (s *server) handleSubmitPoll() http.HandlerFunc {
	type request struct {
		Poll []struct {
			Word string `json:"word"`
			Mark int    `json:"mark"`
		} `json:"poll"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.logger.Error(w, r, http.StatusBadRequest, err)
			return
		}

		p := model.Poll{
			Answer: make([]model.Answer, 0),
		}
		for _, a := range req.Poll {
			pa := model.Answer{
				Word: a.Word,
				Mark: a.Mark,
			}
			p.Answer = append(p.Answer, pa)
		}

		acc := r.Context().Value(contextKeyAccount).(*model.Account)
		for _, pa := range p.Answer {
			err := s.ustorage.Poll().Submit(acc.UUID, &pa)
			if err != nil {
				s.logger.Error(w, r, http.StatusBadRequest, err)
				return
			}
		}

	}
}

func (s *server) handleEducationPlaceGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		earr, err := s.ustorage.EducationPlace().Get()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusCreated, earr)
	}
}

func (s *server) handleAllAccountsGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		aarr, err := s.mstorage.Account().GetAll()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusCreated, aarr)
	}
}

func (s *server) handleStatsGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sarr, err := s.astorage.Poll().GetPollStats()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		// s.logger.Info(earr)
		s.respond(w, r, http.StatusCreated, sarr)
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
