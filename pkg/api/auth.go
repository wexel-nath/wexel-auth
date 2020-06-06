package api

import (
	"net/http"

	"wexel-auth/pkg/jwt"
	"wexel-auth/pkg/logger"
	"wexel-auth/pkg/permission"
	"wexel-auth/pkg/session"
	"wexel-auth/pkg/user"

	"github.com/wexel-nath/auth"
)

type authResponse struct {
	User    auth.User `json:"user"`
	Jwt     string    `json:"jwt"`
	Refresh string    `json:"refresh_token"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func doLogin(request loginRequest) (authResponse, error) {
	userModel, err := user.Authenticate(request.Username, request.Password)
	if err != nil {
		logger.Warn(err)
		return authResponse{}, err
	}

	permissions, err := permission.GetAllForUser(userModel.UserID)
	if err != nil {
		return authResponse{}, err
	}

	jwtUser := buildJwtUser(userModel, permissions)
	jwtToken, err := jwt.Sign(jwtUser)
	if err != nil {
		return authResponse{}, err
	}

	userSession, err := session.Create(userModel.UserID)
	if err != nil {
		return authResponse{}, err
	}

	tokens := authResponse{
		User:    jwtUser,
		Jwt:     jwtToken,
		Refresh: userSession.SessionID,
	}
	return tokens, nil
}

func login(w http.ResponseWriter, r *http.Request) {
	var request loginRequest
	err := unmarshalRequestBody(r, &request)
	if err != nil {
		logger.Error(err)
		jsonResponse(w, http.StatusNotAcceptable, nil, http.StatusText(http.StatusNotAcceptable))
		return
	}

	response, err := doLogin(request)
	if err != nil {
		status := http.StatusUnauthorized
		if err != user.ErrInvalidDetails {
			status = http.StatusInternalServerError
			logger.Error(err)
		}

		jsonResponse(w, status, nil, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, response, "")
}

func doRefresh(jwtUser auth.User, refreshToken string) (authResponse, error) {
	s, err := session.ExtendCurrentSession(refreshToken, jwtUser.UserID)
	if err != nil {
		return authResponse{}, err
	}

	jwtToken, err := jwt.Sign(jwtUser)
	if err != nil {
		return authResponse{}, err
	}

	tokens := authResponse{
		Jwt:     jwtToken,
		Refresh: s.SessionID,
	}
	return tokens, nil
}

func refresh(w http.ResponseWriter, r *http.Request) {
	jwtUser, err := jwt.Authenticate(r)
	if err != nil && err != jwt.ErrExpiredToken {
		logger.Error(err)
		jsonResponse(w, http.StatusUnauthorized, nil, err.Error())
		return
	}

	var request struct {
		RefreshToken string `json:"refresh_token"`
	}
	err = unmarshalRequestBody(r, &request)
	if err != nil {
		logger.Error(err)
		jsonResponse(w, http.StatusNotAcceptable, nil, err.Error())
		return
	}

	response, err := doRefresh(jwtUser, request.RefreshToken)
	if err != nil {
		status := http.StatusUnauthorized
		if err != session.ErrSessionExpired {
			status = http.StatusInternalServerError
			logger.Error(err)
		}

		jsonResponse(w, status, nil, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, response, "")
}

func logout(w http.ResponseWriter, r *http.Request) {
	jwtUser, err := jwt.Authenticate(r)
	if err != nil && err != jwt.ErrExpiredToken {
		logger.Error(err)
		jsonResponse(w, http.StatusUnauthorized, nil, err.Error())
		return
	}

	var request struct {
		RefreshToken string `json:"refresh_token"`
	}
	err = unmarshalRequestBody(r, &request)
	if err != nil {
		logger.Error(err)
		jsonResponse(w, http.StatusNotAcceptable, nil, err.Error())
		return
	}

	_, err = session.EndCurrentSession(request.RefreshToken, jwtUser.UserID)
	if err != nil {
		logger.Error(err)
		jsonResponse(w, http.StatusUnauthorized, nil, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, nil, "")
}

func buildJwtUser(userModel user.User, permissions permission.UserPermissions) auth.User {
	return auth.User{
		UserID:      userModel.UserID,
		FirstName:   userModel.FirstName,
		LastName:    userModel.LastName,
		Email:       userModel.Email,
		Username:    userModel.Username,
		Permissions: auth.Permissions(permissions),
	}
}
