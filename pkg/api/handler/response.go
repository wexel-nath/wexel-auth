package handler

import (
	"encoding/json"
	"net/http"

	"github.com/wexel-nath/wexel-auth/pkg/logger"
)

type authResponse struct {
	Jwt     string `json:"jwt_token"`
	Refresh string `json:"refresh_token"`
}

type response struct {
	Result   interface{} `json:"result"`
	Messages []string    `json:"messages"`
}

func newResponse(result interface{}, messages []string) response {
	return response{
		Result:   result,
		Messages: messages,
	}
}

func writeJsonResponse(
	resp http.ResponseWriter,
	result interface{},
	messages []string,
	status int,
) {
	bytes, err := json.Marshal(newResponse(result, messages))
	if err != nil {
		logger.Error(err)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(status)
	resp.Write(bytes)
}
