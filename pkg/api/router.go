package api

import (
	"log"
	"net/http"

	"github.com/wexel-nath/authrouter"
	"github.com/wexel-nath/wexel-auth/pkg/config"
	"github.com/wexel-nath/wexel-auth/pkg/logger"
)

var (
	routes = []authrouter.Route{
		{
			Method:  http.MethodGet,
			Path:    "/healthz",
			Handler: healthz,
		},
		{
			Method:  http.MethodPost,
			Path:    "/login",
			Handler: login,
		},
		{
			Method:  http.MethodPost,
			Path:    "/refresh",
			Handler: refresh,
		},
	}

	authenticatedRoutes = []authrouter.Route{
		{
			Method:  http.MethodPost,
			Path:    "/logout",
			Handler: logout,
		},
		{
			Method:  http.MethodGet,
			Path:    "/user",
			Handler: getUser,
		},
	}

	authorizedRoutes = []authrouter.Route{
		{
			Method:     http.MethodPost,
			Path:       "/user",
			Service:    "",
			Capability: "user.create",
			Handler:    createUser,
		},
	}
)

func GetRouter() *authrouter.Router {
	auth, err := authrouter.NewAuthenticator(config.GetPublicKeyPath())
	if err != nil {
		log.Fatal(err)
	}

	routerConfig := authrouter.Config{
		Routes:              routes,
		AuthenticatedRoutes: authenticatedRoutes,
		AuthorizedRoutes:    authorizedRoutes,
		EnableCors:          true,
	}

	return authrouter.New(auth, logger.Logger{}, routerConfig)
}
