package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/wexel-nath/wexel-auth/pkg/jwt"
	"github.com/wexel-nath/wexel-auth/pkg/logger"
	"github.com/wexel-nath/wexel-auth/pkg/permission"
	"github.com/wexel-nath/wexel-auth/pkg/session"
	"github.com/wexel-nath/wexel-auth/pkg/user"
)

type authResponse struct {
	User    jwt.User `json:"user"`
	Jwt     string    `json:"jwt"`
	Refresh string    `json:"refresh_token"`
}

func login(r *http.Request) (interface{}, interface{}, int) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		logger.Error(err)
		return nil, err.Error(), http.StatusBadRequest
	}

	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err = json.Unmarshal(body, &request)
	if err != nil {
		logger.Error(err)
		return nil, err.Error(), http.StatusBadRequest
	}

	userModel, err := user.Authenticate(request.Username, request.Password)
	if err != nil {
		logger.Error(err)
		return nil, err.Error(), http.StatusUnauthorized
	}

	permissions, err := permission.GetAllForUser(userModel.UserID)
	if err != nil {
		logger.Error(err)
		return nil, err.Error(), http.StatusInternalServerError
	}

	jwtUser := buildJwtUser(userModel, permissions)
	jwt, err := jwt.Sign(jwtUser)
	if err != nil {
		logger.Error(err)
		return nil, err.Error(), http.StatusInternalServerError
	}

	userSession, err := session.Create(userModel.UserID)
	if err != nil {
		logger.Error(err)
		return nil, err.Error(), http.StatusInternalServerError
	}

	tokens := authResponse{
		User:    jwtUser,
		Jwt:     jwt,
		Refresh: userSession.SessionID,
	}
	return tokens, nil, http.StatusOK
}

func refresh(r *http.Request) (interface{}, interface{}, int) {
	jwtUser, err := jwt.Authenticate(r)
	if err != nil && err != jwt.ErrExpiredToken {
		logger.Error(err)
		return nil, err.Error(), http.StatusUnauthorized
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		logger.Error(err)
		return nil, err.Error(), http.StatusBadRequest
	}

	var request struct {
		RefreshToken string `json:"refresh_token"`
	}
	err = json.Unmarshal(body, &request)
	if err != nil {
		logger.Error(err)
		return nil, err.Error(), http.StatusBadRequest
	}

	s, err := session.ExtendCurrentSession(request.RefreshToken, jwtUser.UserID)
	if err != nil {
		logger.Error(err)
		return nil, err.Error(), http.StatusUnauthorized
	}

	jwt, err := jwt.Sign(jwtUser)
	if err != nil {
		logger.Error(err)
		return nil, err.Error(), http.StatusInternalServerError
	}

	tokens := authResponse{
		Jwt:     jwt,
		Refresh: s.SessionID,
	}
	return tokens, nil, http.StatusOK
}

func logout(r *http.Request, userInterface interface{}) (interface{}, interface{}, int) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		logger.Error(err)
		return nil, err.Error(), http.StatusBadRequest
	}

	var request struct {
		RefreshToken string `json:"refresh_token"`
	}
	err = json.Unmarshal(body, &request)
	if err != nil {
		logger.Error(err)
		return nil, err.Error(), http.StatusBadRequest
	}

	user := userInterface.(jwt.User)
	_, err = session.EndCurrentSession(request.RefreshToken, user.UserID)
	if err != nil {
		logger.Error(err)
		return nil, err.Error(), http.StatusUnauthorized
	}

	return nil, nil, http.StatusOK
}

func buildJwtUser(userModel user.User, permissions permission.UserPermissions) jwt.User {
	return jwt.User{
		UserID:      userModel.UserID,
		FirstName:   userModel.FirstName,
		LastName:    userModel.LastName,
		Email:       userModel.Email,
		Username:    userModel.Username,
		Permissions: permissions,
	}
}
