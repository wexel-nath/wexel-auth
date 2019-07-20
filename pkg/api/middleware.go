package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/wexel-nath/authrouter"
	"github.com/wexel-nath/wexel-auth/pkg/logger"
)

func middlewareWrapper(handler authrouter.Handler) authrouter.Handler {
	next := loggerWrapper(handler)
	return next
}

func loggerWrapper(next authrouter.Handler) authrouter.Handler {
	return func(r *http.Request) (interface{}, interface{}, int) {
		logger.Info("%s %s", r.Method, r.URL.Path)
		return next(r)
	}
}

func corsMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		next(w, r, p)
	}
}
