package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/wexel-nath/wexel-auth/pkg/auth"
	"github.com/wexel-nath/wexel-auth/pkg/logger"
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

type handleFunc func(r *http.Request) (interface{}, int, error)

func requestHandler(handler handleFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		messages := make([]string, 0)

		result, statusCode, err := handler(r)
		if err != nil {
			logger.Error(err)
			messages = []string{ err.Error() }
		}

		writeJsonResponse(w, result, messages, statusCode)
	}
}

type handleUserFunc func(r *http.Request, user auth.User) (interface{}, int, error)

func authRequestHandler(handler handleUserFunc, service string, capability string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var result interface{}
		var statusCode int
		messages := make([]string, 0)

		user, err := auth.GetAuthorizedUser(r, service, capability)
		if err == auth.ErrNotAuthorized {
			statusCode = http.StatusForbidden
		} else if err != nil {
			statusCode = http.StatusUnauthorized
		} else {
			result, statusCode, err = handler(r, user)
		}

		if err != nil {
			logger.Error(err)
			messages = []string{ err.Error() }
		}

		writeJsonResponse(w, result, messages, statusCode)
	}
}
