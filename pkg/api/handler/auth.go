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
	"github.com/wexel-nath/wexel-auth/pkg/session"
	"github.com/wexel-nath/wexel-auth/pkg/user"
)

func HandleLogin(r *http.Request) (interface{}, int, error) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err = json.Unmarshal(body, &request)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	userModel, err := user.Authenticate(request.Username, request.Password)
	if err != nil {
		return nil, http.StatusUnauthorized, err
	}

	jwt, err := auth.SignUser(userModel)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	userSession, err := session.Create(userModel.UserID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	tokens := authResponse{
		Jwt:     jwt,
		Refresh: userSession.SessionID,
	}
	return tokens, http.StatusOK, nil
}

func RefreshHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userModel, err := getAuthenticatedUser(r)
	if err != nil {
		logger.Error(err)
		messages := []string{ err.Error() }
		writeJsonResponse(w, nil, messages, http.StatusUnauthorized)
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

	// todo: check if session is valid

	jwt, err := auth.SignUser(userModel)
	if err != nil {
		logger.Error(err)
		messages := []string { err.Error() }
		writeJsonResponse(w, nil, messages, http.StatusInternalServerError)
		return
	}

	// todo: extend session and return refresh token
	tokens := authResponse{
		Jwt:     jwt,
		Refresh: "TO.DO.INCOMPLETE",
	}
	writeJsonResponse(w, tokens, nil, http.StatusNotImplemented)
}

func getAuthenticatedUser(r *http.Request) (user.User, error) {
	value := r.Header.Get("Authorization")
	s := strings.Split(value, "Bearer")
	if len(s) != 2 {
		return user.User{}, errors.New("authorization header is not of the form: Bearer <token>")
	}

	token := strings.TrimSpace(s[1])
	return auth.Verify(token)
}
