package api

import (
	"encoding/json"
	"net/http"

	"github.com/wexel-nath/wexel-auth/pkg/logger"
)

type response struct {
	Result  interface{} `json:"result"`
	Message string      `json:"message"`
}

func jsonResponse(w http.ResponseWriter, status int, result interface{}, message string) {
	r := response{
		Result:  result,
		Message: message,
	}

	resp, err := json.Marshal(r)
	if err != nil {
		logger.Error(err)

		status = http.StatusInternalServerError
		resp = []byte(`{}`)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(resp)
}
