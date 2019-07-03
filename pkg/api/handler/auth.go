package handler

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/wexel-nath/wexel-auth/pkg/auth"
	"github.com/wexel-nath/wexel-auth/pkg/logger"
	"github.com/wexel-nath/wexel-auth/pkg/user"
)


func LoginHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		logger.Error(err)
		messages := []string { err.Error() }
		writeJsonResponse(w, nil, messages, http.StatusBadRequest)
		return
	}

	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err = json.Unmarshal(body, &request)
	if err != nil {
		logger.Error(err)
		messages := []string { err.Error() }
		writeJsonResponse(w, nil, messages, http.StatusBadRequest)
		return
	}

	userModel, err := user.Authenticate(request.Username, request.Password)
	if err != nil {
		logger.Error(err)
		messages := []string { err.Error() }
		writeJsonResponse(w, nil, messages, http.StatusUnauthorized)
		return
	}

	jwt, err := auth.SignUser(userModel)
	if err != nil {
		logger.Error(err)
		messages := []string { err.Error() }
		writeJsonResponse(w, nil, messages, http.StatusBadRequest)
		return
	}

	// todo: return valid JWT and refresh token
	tokens := authResponse{
		Jwt:     jwt,
		Refresh: "TO.DO.INCOMPLETE",
	}

	err = auth.Verify(jwt)
	if err != nil {
		logger.Error(err)
	}

	writeJsonResponse(w, tokens, nil, http.StatusOK)
}

func RefreshHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ok, err := isAuthenticated(r)
	if !ok || err != nil {
		message := ""
		if err != nil {
			logger.Error(err)
			message = err.Error()
		} else {
			message = "not allowed"
		}
		writeJsonResponse(w, nil, []string{ message }, http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		logger.Error(err)
		messages := []string{ err.Error() }
		writeJsonResponse(w, nil, messages, http.StatusBadRequest)
		return
	}

	var request struct {
		RefreshToken string `json:"refresh_token"`
	}
	err = json.Unmarshal(body, &request)
	if err != nil {
		logger.Error(err)
		messages := []string{ err.Error() }
		writeJsonResponse(w, nil, messages, http.StatusBadRequest)
		return
	}

	// todo: return valid JWT and refresh token
	tokens := authResponse{
		Jwt:     "TO.DO.INCOMPLETE",
		Refresh: "TO.DO.INCOMPLETE",
	}
	writeJsonResponse(w, tokens, nil, http.StatusNotImplemented)
}

func isAuthenticated(r *http.Request) (bool, error) {
	value := r.Header.Get("Authorization")
	s := strings.Split(value, "Bearer")
	if len(s) != 2 {
		return false, errors.New("authorization header is not of the form: Bearer <token>")
	}

	token := strings.TrimSpace(s[1])

	// validate JWT
	logger.Debug("token: %s", token)

	return true, nil
}
