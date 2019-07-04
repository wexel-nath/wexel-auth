package handler

import (
	"net/http"
)

func HandleHealthz(r *http.Request) (interface{}, int, error) {
	result := struct{
		Status string `json:"status"`
	}{
		Status: "ok",
	}
	return result, http.StatusOK, nil
}
