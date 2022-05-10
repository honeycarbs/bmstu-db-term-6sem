package poll

import (
	"api/pkg/logger"
	"api/pkg/middleware/jwt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	FORM_ENDPOINT = "/api/form"
)

// type Form struct {
// 	ID string `json:"id"`
// }

type Handler struct {
	Logger logger.Logger
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, FORM_ENDPOINT, jwt.JWTMiddleware(h.GetPoll))
	router.HandlerFunc(http.MethodPost, FORM_ENDPOINT, jwt.JWTMiddleware(h.SubmitPoll))
}

func (h *Handler) GetPoll(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Get form")
	if r.Context().Value("user_id") == nil {
		h.Logger.Error("there is no user_uuid in context")
		return
	}
	h.Logger.Infof("User id: %s", r.Context().Value("user_id").(string))
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("i will put form there"))
	w.WriteHeader(200)
}

func (h *Handler) SubmitPoll(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("This form will be updated when user clicks the answer")
	if r.Context().Value("user_id") == nil {
		h.Logger.Error("there is no user_uuid in context")
		return
	}
	h.Logger.Infof("User id: %s", r.Context().Value("user_id").(string))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
}
