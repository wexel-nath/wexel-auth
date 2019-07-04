package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/wexel-nath/wexel-auth/pkg/logger"
	"net/http"
)

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

func requestHandler(handler func(r *http.Request) (interface{}, int, error)) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		messages := []string(nil)

		result, statusCode, err := handler(r)
		if err != nil {
			logger.Error(err)
			messages = []string{ err.Error() }
		}

		writeJsonResponse(w, result, messages, statusCode)
	}
}
