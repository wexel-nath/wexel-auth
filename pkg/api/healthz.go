package api

import (
	"net/http"
)

func healthz(r *http.Request) (interface{}, interface{}, int) {
	result := struct{
		Status string `json:"status"`
	}{
		Status: "ok",
	}
	return result, nil, http.StatusOK
}
