package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/wexel-nath/wexel-auth/pkg/api"
	"github.com/wexel-nath/wexel-auth/pkg/auth"
	"github.com/wexel-nath/wexel-auth/pkg/config"
	"github.com/wexel-nath/wexel-auth/pkg/session"
)

func main() {
	config.Configure()
	auth.Configure()
	session.Configure()

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
	fmt.Println("Listening on " + address)
	log.Fatal(http.ListenAndServe(address, api.GetRouter()))
}
