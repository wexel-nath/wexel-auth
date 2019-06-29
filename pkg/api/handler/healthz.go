package handler

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/wexel-nath/wexel-auth/pkg/logger"
)

func HealthzHandler(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write([]byte(`{"status":"ok"}`))
	if err != nil {
		logger.Error(err)
	}
}
