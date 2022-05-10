package auth

import (
	"api/internal/config"
	"api/pkg/cache"
	"api/pkg/logger"
	"encoding/json"
	"net/http"
	"time"

	myjwt "api/pkg/middleware/jwt"

	"github.com/cristalhq/jwt/v3"
	"github.com/julienschmidt/httprouter"
)

const (
	AUTH_ENDPOINT = "/api/auth"
)

type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Handler struct {
	Logger  logger.Logger
	RTCache cache.Repository
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, AUTH_ENDPOINT, h.Auth)
}

func (h *Handler) Auth(w http.ResponseWriter, r *http.Request) {
	var u user
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		h.Logger.Fatal(err)
	}
	defer r.Body.Close()
	if u.Username != "honeycarbs" && u.Password != "honeycarbs" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	key := []byte(config.GetConfig().JWT.Secret)
	signer, err := jwt.NewSignerHS(jwt.HS256, key)
	if err != nil {
		w.WriteHeader(418)
		return
	}

	builder := jwt.NewBuilder(signer)

	claims := myjwt.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        "uuid_here",
			Audience:  []string{"users"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * 5)),
		},
		Email: "example@email.here",
	}

	token, err := builder.Build(claims)
	if err != nil {
		h.Logger.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	jsonBytes, err := json.Marshal(map[string]string{
		"token": token.String(),
	})
	if err != nil {
		h.Logger.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
	w.WriteHeader(200)
}
