package api

import (
	"net/http"
)

func healthz(w http.ResponseWriter, r *http.Request) {
	result := struct{
		Status string `json:"status"`
	}{
		Status: "ok",
	}

	jsonResponse(w, http.StatusOK, result, "")
}
