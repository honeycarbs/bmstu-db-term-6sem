package jwt

import (
	"api/internal/config"
	"api/pkg/logger"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/cristalhq/jwt/v3"
)

type UserClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
}

func JWTMiddleware(h http.HandlerFunc) http.HandlerFunc {
	logger := logger.GetLogger()
	return func(w http.ResponseWriter, r *http.Request) {
		authHandler := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		logger.Info(authHandler)
		if len(authHandler) != 2 {
			logger.Error("Malformed token")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("malformed token"))
			return
		}

		logger.Debug("create JWT verifier")
		jwtToken := authHandler[1]
		key := []byte(config.GetConfig().JWT.Secret)
		verifier, err := jwt.NewVerifierHS(jwt.HS256, key)
		if err != nil {
			unautorized(w, err)
		}

		logger.Debug("parce and verify JWT token")
		token, err := jwt.ParseAndVerifyString(jwtToken, verifier)
		if err != nil {
			unautorized(w, err)
			return
		}

		var uc UserClaims
		err = json.Unmarshal(token.RawClaims(), &uc)
		if err != nil {
			unautorized(w, err)
			return
		}

		if valid := uc.IsValidAt(time.Now()); !valid {
			logger.Error("Malformed token")
			unautorized(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", uc.ID)
		// !!! put role in context (decompose)
		// put role in context eto tozhe norm
		// or get user from db and check role
		h(w, r.WithContext(ctx))
	}
}

func unautorized(w http.ResponseWriter, err error) {
	logger.GetLogger().Error(err)
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("unauthorized"))
}
