package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/wexel-nath/wexel-auth/pkg/logger"
)

func middleware(handler httprouter.Handle) httprouter.Handle {
	next := loggerWrapper(handler)
	next = corsMiddleware(next)
	return next
}

func loggerWrapper(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		logger.Info("%s %s", r.Method, r.URL.Path)
		next(w, r, p)
	}
}

func corsMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		next(w, r, p)
	}
}
