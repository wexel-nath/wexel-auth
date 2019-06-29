package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/wexel-nath/wexel-auth/pkg/api"
	"github.com/wexel-nath/wexel-auth/pkg/config"
)

func main() {
	config.Configure()

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
