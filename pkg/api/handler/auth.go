package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/wexel-nath/wexel-auth/pkg/auth"
	"github.com/wexel-nath/wexel-auth/pkg/permission"
	"github.com/wexel-nath/wexel-auth/pkg/session"
	"github.com/wexel-nath/wexel-auth/pkg/user"
)

type authResponse struct {
	User    auth.User `json:"user"`
	Jwt     string    `json:"jwt"`
	Refresh string    `json:"refresh_token"`
}

func Login(r *http.Request) (interface{}, int, error) {
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

	permissions, err := permission.GetAllForUser(userModel.UserID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	authUser := buildAuthUserModel(userModel, permissions)
	jwt, err := auth.Sign(authUser)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	userSession, err := session.Create(userModel.UserID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	tokens := authResponse{
		User:    authUser,
		Jwt:     jwt,
		Refresh: userSession.SessionID,
	}
	return tokens, http.StatusOK, nil
}

func Refresh(r *http.Request) (interface{}, int, error) {
	authUser, err := auth.GetAuthenticatedUser(r)
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

	s, err := session.ExtendCurrentSession(request.RefreshToken, authUser.UserID)
	if err != nil {
		return nil, http.StatusUnauthorized, err
	}

	jwt, err := auth.Sign(authUser)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	tokens := authResponse{
		Jwt:     jwt,
		Refresh: s.SessionID,
	}
	return tokens, http.StatusOK, nil
}

func Logout(r *http.Request, user auth.User) (interface{}, int, error) {
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

	_, err = session.EndCurrentSession(request.RefreshToken, user.UserID)
	if err != nil {
		return nil, http.StatusUnauthorized, err
	}

	return nil, http.StatusOK, nil
}

func buildAuthUserModel(userModel user.User, permissions permission.UserPermissions) auth.User {
	return auth.User{
		UserID:      userModel.UserID,
		FirstName:   userModel.FirstName,
		LastName:    userModel.LastName,
		Email:       userModel.Email,
		Username:    userModel.Username,
		Permissions: permissions,
	}
}
