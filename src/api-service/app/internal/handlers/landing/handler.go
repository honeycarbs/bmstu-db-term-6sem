package landing

import (
	"api/pkg/logger"
	"api/pkg/middleware/jwt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	LANDING_ENDPOINT = "/api/landing"
)

type Handler struct {
	Logger logger.Logger
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, LANDING_ENDPOINT, jwt.JWTMiddleware(h.GetLanding))
	// router.HandlerFunc(http.MethodPost, FORM_ENDPOINT, jwt.JWTMiddleware(h.SubmitForm))
}

func (h *Handler) GetLanding(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Landing")
	if r.Context().Value("user_id") == nil {
		h.Logger.Error("there is no user_uuid in context")
		return
	}
	h.Logger.Infof("User id: %s", r.Context().Value("user_id").(string))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
}
