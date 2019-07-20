package api

import (
	"net/http"

	"github.com/wexel-nath/authrouter"
)

func healthz(_ *http.Request, _ authrouter.User) (interface{}, interface{}, int) {
	result := struct{
		Status string `json:"status"`
	}{
		Status: "ok",
	}
	return result, nil, http.StatusOK
}
