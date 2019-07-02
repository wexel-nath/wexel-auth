package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
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

	// should not return userModel
	// todo: return a JWT and refresh token

	writeJsonResponse(w, userModel, nil, http.StatusOK)
}

func RefreshHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// todo: capture and authenticate JWT from Header

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		logger.Error(err)
		messages := []string { err.Error() }
		writeJsonResponse(w, nil, messages, http.StatusBadRequest)
		return
	}

	var request struct {
		RefreshToken string `json:"refresh_token"`
	}
	err = json.Unmarshal(body, &request)
	if err != nil {
		logger.Error(err)
		messages := []string { err.Error() }
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
