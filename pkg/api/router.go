package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/wexel-nath/wexel-auth/pkg/api/handler"
)

func GetRouter() *httprouter.Router {
	router := httprouter.New()
	endpointMethods := map[string][]string{}

	for _, route := range getRoutes() {
		endpointMethods[route.path] = append(endpointMethods[route.path], route.method)
		router.Handle(route.method, route.path, middlewareWrapper(route.handler))
	}

	for path, methods := range endpointMethods {
		router.OPTIONS(path, constructOptions(methods))
	}

	return router
}

type route struct {
	method  string
	path    string
	handler httprouter.Handle
}

func getRoutes() []route {
	return []route{
		{
			method:  http.MethodGet,
			path:    "/healthz",
			handler: requestHandler(handler.Healthz),
		},
		{
			method: http.MethodPost,
			path:    "/login",
			handler: requestHandler(handler.Login),
		},
		{
			method: http.MethodPost,
			path:    "/refresh",
			handler: requestHandler(handler.Refresh),
		},
		{
			method: http.MethodPost,
			path:    "/logout",
			handler: authRequestHandler(handler.Logout, "", ""),
		},
		{
			method:  http.MethodPost,
			path:    "/user",
			handler: authRequestHandler(handler.CreateUser, "", "user.create"),
		},
		{
			method:  http.MethodGet,
			path:    "/user",
			handler: authRequestHandler(handler.GetUser, "", ""),
		},
	}
}

func constructOptions(methods []string) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	methodCsv := strings.Join(append(methods, "OPTIONS"), ",")
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", methodCsv)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, "{}")
	}
}
