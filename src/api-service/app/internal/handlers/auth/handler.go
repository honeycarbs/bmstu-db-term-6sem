package auth

import (
	"api/internal/config"
	"api/pkg/cache"
	"api/pkg/logger"
	myjwt "api/pkg/middleware/jwt"
	"encoding/json"
	"net/http"
	"time"

	"github.com/cristalhq/jwt/v3"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

const (
	AUTH_ENDPOINT   = "/api/auth"
	SIGNUP_ENDPOINT = "/api/signup"
)

type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type NewUser struct {
	user
	Email string `json:"email"`
}

type Handler struct {
	Logger  logger.Logger
	RTCache cache.Repository
}

type refresh struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, AUTH_ENDPOINT, h.Auth)
	router.HandlerFunc(http.MethodPut, AUTH_ENDPOINT, h.Auth)
	router.HandlerFunc(http.MethodPost, SIGNUP_ENDPOINT, h.Signup)
}

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	var nu NewUser
	if err := json.NewDecoder(r.Body).Decode(&nu); err != nil {
		h.Logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	jsonBytes, rc := h.generateAccessToken()
	if rc != 0 {
		w.WriteHeader(rc)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(jsonBytes)

}

func (h *Handler) Auth(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var u user
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			h.Logger.Fatal(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		if u.Username != "honeycarbs" && u.Password != "honeycarbs" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

	case http.MethodPut:
		var refreshTokenWrapper refresh
		if err := json.NewDecoder(r.Body).Decode(&refreshTokenWrapper); err != nil {
			h.Logger.Fatal(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		userIdBytes, err := h.RTCache.Get([]byte(refreshTokenWrapper.RefreshToken))
		h.Logger.Info("refresh token user_id: %s", userIdBytes)
		if err != nil {
			h.Logger.Error(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		h.RTCache.Delete([]byte(refreshTokenWrapper.RefreshToken))
	}

	jsonBytes, ret := h.generateAccessToken()
	if ret != 0 {
		w.WriteHeader(ret)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(jsonBytes)

}

func (h *Handler) generateAccessToken() ([]byte, int) {
	key := []byte(config.GetConfig().JWT.Secret)
	signer, err := jwt.NewSignerHS(jwt.HS256, key)
	if err != nil {
		return nil, 418
	}

	builder := jwt.NewBuilder(signer)
	claims := myjwt.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        "uuid_here",
			Audience:  []string{"users"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		Email: "example@email.here",
	}

	token, err := builder.Build(claims)
	if err != nil {
		h.Logger.Error(err)
		return nil, http.StatusUnauthorized
	}

	h.Logger.Info("create refresh token")
	refreshTokenUUID := uuid.New()
	err = h.RTCache.Set([]byte(refreshTokenUUID.String()), []byte("user_uuid"), 0)
	if err != nil {
		h.Logger.Error(err)
		return nil, http.StatusInternalServerError
	}

	jsonBytes, err := json.Marshal(map[string]string{
		"token":         token.String(),
		"refresh_token": refreshTokenUUID.String(),
	})

	if err != nil {
		return nil, http.StatusInternalServerError
	}

	return jsonBytes, 0
}
