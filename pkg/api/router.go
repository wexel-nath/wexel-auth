package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/wexel-nath/authrouter"
	"github.com/wexel-nath/jwt"
)

func GetRouter(publicKeyPath string) *authrouter.Router {
	auth, err := jwt.NewAuthenticator(publicKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	router := authrouter.NewRouter(auth)
	endpointMethods := map[string][]string{}

	for authMethod, routes := range getRoutes() {
		for _, route := range routes {
			endpointMethods[route.path] = append(endpointMethods[route.path], route.method)

			switch authMethod {
			case "handle":
				router.Handle(route.method, route.path, route.handler)
			case "authenticated":
				router.HandleWithAuthentication(route.method, route.path, route.uHandler)
			case "authorized":
				router.HandleWithAuthorization(
					route.method,
					route.path,
					route.uHandler,
					route.service,
					route.capability,
				)
			}
		}
	}

	for path, methods := range endpointMethods {
		router.HttpRouter.OPTIONS(path, constructOptions(methods))
	}

	return router
}

type route struct {
	method     string
	path       string
	handler    authrouter.Handler
	service    string
	capability string
	uHandler   authrouter.HandlerWithUser
}

func getRoutes() map[string][]route {
	return map[string][]route{
		"handle": {
			{
				method:  http.MethodGet,
				path:    "/healthz",
				handler: healthz,
			},
			{
				method:  http.MethodPost,
				path:    "/login",
				handler: login,
			},
			{
				method:  http.MethodPost,
				path:    "/refresh",
				handler: refresh,
			},
		},
		"authenticated": {
			{
				method:   http.MethodPost,
				path:     "/logout",
				uHandler: logout,
			},
			{
				method:   http.MethodGet,
				path:     "/user",
				uHandler: getUser,
			},
		},
		"authorized": {
			{
				method:     http.MethodPost,
				path:       "/user",
				service:    "",
				capability: "user.create",
				uHandler:   createUser,
			},
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
