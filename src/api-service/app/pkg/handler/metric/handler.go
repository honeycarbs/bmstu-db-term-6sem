package metric

import (
	"api/pkg/logger"
	"api/pkg/middleware/jwt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	ENDPOINT = "/api/heartbeat"
)

type Handler struct {
	Logger logger.Logger
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, ENDPOINT, jwt.JWTMiddleware(h.Heartbeat))
}

func (h *Handler) Heartbeat(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(204)
}

// func Test(w http.ResponseWriter, req *http.Request) {
// 	w.WriteHeader(204)
// }
