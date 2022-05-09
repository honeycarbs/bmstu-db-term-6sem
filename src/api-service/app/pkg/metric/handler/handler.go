package metric

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	HEARTBEET_URL = "/api/heartbeet"
)

func Heartbeet(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	w.WriteHeader(204)
}
