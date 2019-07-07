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
		router.Handle(http.MethodOptions, route.pattern, enableCors)
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
			handler: requestHandler(handler.HandleHealthz),
		},
		{
			method: http.MethodPost,
			pattern: "/login",
			handler: requestHandler(handler.HandleLogin),
		},
		{
			method: http.MethodPost,
			pattern: "/refresh",
			handler: requestHandler(handler.HandleRefresh),
		},
		{
			method: http.MethodPost,
			pattern: "/user",
			handler: authRequestHandler(handler.HandleCreateUser, "", "user.create"),
		},
	}
}

func enableCors(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.WriteHeader(http.StatusNoContent)
}
