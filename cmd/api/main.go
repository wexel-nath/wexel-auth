package main

import (
	"log"
	"net/http"

	"wexel-auth/pkg/api"
	"wexel-auth/pkg/config"
	"wexel-auth/pkg/jwt"
	"wexel-auth/pkg/logger"
	"wexel-auth/pkg/session"
)

func main() {
	config.Configure()
	session.Configure()
	jwt.Configure()

	startServer()
}

func getListenAddress() string {
	port := config.GetPort()

	if len(port) == 0 {
		log.Fatal("PORT must be set")
	}

	return ":" + port
}

func startServer() {
	address := getListenAddress()
	logger.Info("Listening on %s", address)
	log.Fatal(http.ListenAndServe(address, api.GetRouter()))
}
