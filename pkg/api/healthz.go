package api

import (
	"net/http"

	"wexel-auth/pkg/config"
)

func healthz(w http.ResponseWriter, r *http.Request) {
	result := struct{
		Status string `json:"status"`
		Image  string `json:"image"`
	}{
		Status: "ok",
		Image:  config.GetImageTag(),
	}

	jsonResponse(w, http.StatusOK, result, "")
}
