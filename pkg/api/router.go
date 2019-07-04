package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/wexel-nath/wexel-auth/pkg/api/handler"
)

func GetRouter() *httprouter.Router {
	router := httprouter.New()

	for _, route := range getRoutes() {
		handle := requestHandler(route.handler)
		router.Handle(route.method, route.pattern, middleware(handle))
	}

	return router
}

type route struct {
	method  string
	pattern string
	handler func(r *http.Request) (interface{}, int, error)
}

func getRoutes() []route {
	return []route{
		{
			method:  http.MethodGet,
			pattern: "/healthz",
			handler: handler.HandleHealthz,
		},
		{
			method: http.MethodPost,
			pattern: "/login",
			handler: handler.HandleLogin,
		},
		{
			method: http.MethodPost,
			pattern: "/refresh",
			handler: handler.HandleRefresh,
		},
		{
			method: http.MethodPost,
			pattern: "/user",
			handler: handler.HandleCreateUser,
		},
	}
}
