package api

import (
	"net/http"
	"strings"

	"github.com/wexel-nath/authrouter"
	"github.com/wexel-nath/wexel-auth/pkg/logger"
	"github.com/wexel-nath/wexel-auth/pkg/permission"
	"github.com/wexel-nath/wexel-auth/pkg/user"
)

func createUser(r *http.Request, _ authrouter.User) (interface{}, interface{}, int) {
	var request struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Username  string `json:"username"`
		Service   string `json:"service_name"`
	}
	err := unmarshalRequestBody(r, &request)
	if err != nil {
		logger.Error(err)
		return nil, err.Error(), http.StatusBadRequest
	}

	// todo: this should generate a random password instead, after forgot/change password is ready
	password := "4Me2Change"

	userModel, err := user.Create(
		request.FirstName,
		request.LastName,
		request.Email,
		request.Username,
		password,
	)
	if err != nil {
		logger.Error(err)
		return nil, err.Error(), http.StatusBadRequest
	}

	err = permission.AddUserPermission(userModel.UserID, strings.ToLower(request.Service))
	if err != nil {
		logger.Error(err)
	}

	return userModel, nil, http.StatusCreated
}

func getUser(_ *http.Request, authUser authrouter.User) (interface{}, interface{}, int) {
	return authUser, nil, http.StatusOK
}

func changePassword(r *http.Request, authUser authrouter.User) (interface{}, interface{}, int) {
	var request struct {
		NewPassword string `json:"new_password"`
	}
	err := unmarshalRequestBody(r, &request)
	if err != nil {
		logger.Error(err)
		return nil, err.Error(), http.StatusBadRequest
	}

	err = user.ChangePassword(authUser.UserID, request.NewPassword)
	if err != nil {
		logger.Error(err)
		return nil, err.Error(), http.StatusBadRequest
	}

	return authUser, nil, http.StatusOK
}
