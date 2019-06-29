package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/wexel-nath/wexel-auth/pkg/api/handler"
)

func GetRouter() *httprouter.Router {
	router := httprouter.New()

	for _, route := range getRoutes() {
		router.Handle(route.method, route.pattern, middleware(route.handler))
	}

	return router
}

type route struct {
	method  string
	pattern string
	handler httprouter.Handle
}

func getRoutes() []route {
	return []route{
		{
			method:  http.MethodGet,
			pattern: "/healthz",
			handler: handler.HealthzHandler,
		},
	}
}
