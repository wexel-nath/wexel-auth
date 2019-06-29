package api

import (
	"github.com/julienschmidt/httprouter"
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
	return []route{}
}
