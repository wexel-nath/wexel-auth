package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/wexel-nath/wexel-auth/pkg/logger"
)

func middleware(handler httprouter.Handle) httprouter.Handle {
	next := loggerWrapper(handler)
	return next
}

func loggerWrapper(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		logger.Info("%s %s", r.Method, r.URL.Path)
		next(w, r, p)
	}
}
