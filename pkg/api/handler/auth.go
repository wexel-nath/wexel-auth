package handler

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/wexel-nath/wexel-auth/pkg/auth"
	"github.com/wexel-nath/wexel-auth/pkg/session"
	"github.com/wexel-nath/wexel-auth/pkg/user"
)

type authResponse struct {
	Jwt     string `json:"jwt_token"`
	Refresh string `json:"refresh_token"`
}

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

func HandleRefresh(r *http.Request) (interface{}, int, error) {
	userModel, err := getAuthenticatedUser(r)
	if err != nil && err != auth.ErrExpiredToken {
		return nil, http.StatusUnauthorized, err
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	var request struct {
		RefreshToken string `json:"refresh_token"`
	}
	err = json.Unmarshal(body, &request)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	s, err := session.ExtendCurrentSession(request.RefreshToken, userModel.UserID)
	if err != nil {
		return nil, http.StatusUnauthorized, err
	}

	jwt, err := auth.SignUser(userModel)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	tokens := authResponse{
		Jwt:     jwt,
		Refresh: s.SessionID,
	}
	return tokens, http.StatusOK, nil
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
