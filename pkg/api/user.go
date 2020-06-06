package api

import (
	"net/http"

	"wexel-auth/pkg/jwt"
	"wexel-auth/pkg/logger"
	"wexel-auth/pkg/permission"
	"wexel-auth/pkg/user"
)

type createUserRequest struct {
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	Email       string  `json:"email"`
	Permissions []int64 `json:"permissions"`
}

func doCreateUser(request createUserRequest) (user.User, error) {
	userModel, err := user.Create(
		request.FirstName,
		request.LastName,
		request.Email,
	)
	if err != nil {
		return user.User{}, err
	}

	err = permission.AddUserPermissions(userModel.UserID, request.Permissions)
	if err != nil {
		logger.Error(err)
	}

	return userModel, nil
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var request createUserRequest
	err := unmarshalRequestBody(r, &request)
	if err != nil {
		logger.Error(err)
		jsonResponse(w, http.StatusNotAcceptable, nil, err.Error())
		return
	}

	newUser, err := doCreateUser(request)
	if err != nil {
		logger.Error(err)
		jsonResponse(w, http.StatusInternalServerError, nil, err.Error())
		return
	}

	jsonResponse(w, http.StatusCreated, newUser, "")
}

func getUser(w http.ResponseWriter, r *http.Request) {
	u, err := jwt.UserFromContext(r.Context())
	if err != nil {
		logger.Error(err)
		jsonResponse(w, http.StatusUnauthorized, nil, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, u, "")
}

func changePassword(w http.ResponseWriter, r *http.Request) {
	u, err := jwt.UserFromContext(r.Context())
	if err != nil {
		logger.Error(err)
		jsonResponse(w, http.StatusUnauthorized, nil, err.Error())
		return
	}

	var request struct {
		NewPassword string `json:"new_password"`
	}
	err = unmarshalRequestBody(r, &request)
	if err != nil {
		logger.Error(err)
		jsonResponse(w, http.StatusNotAcceptable, nil, err.Error())
		return
	}

	err = user.ChangePassword(u.UserID, request.NewPassword)
	if err != nil {
		logger.Error(err)
		jsonResponse(w, http.StatusNotAcceptable, nil, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, u, "")
}

func getAllUsers(w http.ResponseWriter, _ *http.Request) {
	users, err := user.GetAll()
	if err != nil {
		logger.Error(err)
		jsonResponse(w, http.StatusInternalServerError, nil, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, users, "")
}
