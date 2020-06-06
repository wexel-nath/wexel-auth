package api

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"

	"wexel-auth/pkg/jwt"
)

func GetRouter() chi.Router {
	router := chi.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	})
	router.Use(c.Handler)

	addRoutes(router)
	return router
}

func addRoutes(r chi.Router) {
	r.Get("/healthz", healthz)

	r.Post("/login", login)
	r.Post("/refresh", refresh)
	r.Post("/logout", logout)

	// authorized routes
	r.With(jwt.AuthorizationMiddleware("user.create")).
		Get("/service/{serviceName}/permissions", getServicePermissions)

	r.With(jwt.AuthorizationMiddleware("user.create")).
		Post("/user", createUser)

	r.Route("/", func(r chi.Router) {
		r.Use(jwt.Middleware)

		r.Get("/user", getUser)
		r.Get("/users", getAllUsers)
		r.Post("/change-password", changePassword)
	})
}
