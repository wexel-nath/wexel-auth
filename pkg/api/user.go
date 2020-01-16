package api

import (
	"net/http"

	"github.com/wexel-nath/wexel-auth/pkg/jwt"
	"github.com/wexel-nath/wexel-auth/pkg/logger"
	"github.com/wexel-nath/wexel-auth/pkg/permission"
	"github.com/wexel-nath/wexel-auth/pkg/user"
)

type createUserRequest struct {
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	Email       string  `json:"email"`
	Password    string  `json:"password"`
	Username    string  `json:"username"`
	Permissions []int64 `json:"permissions"`
}

func doCreateUser(request createUserRequest) (user.User, error) {
	// todo: this should generate a random password instead, after forgot/change password is ready
	if request.Password == "" {
		request.Password = "4Me2Change"
	}

	userModel, err := user.Create(
		request.FirstName,
		request.LastName,
		request.Email,
		request.Username,
		request.Password,
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
	_, err := jwt.Authorize(r, "user.create")
	if err != nil && err != jwt.ErrExpiredToken {
		if err != jwt.ErrExpiredToken {
			logger.Error(err)
		}
		jsonResponse(w, http.StatusUnauthorized, nil, err.Error())
		return
	}

	var request createUserRequest
	err = unmarshalRequestBody(r, &request)
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
	jwtUser, err := jwt.Authenticate(r)
	if err != nil && err != jwt.ErrExpiredToken {
		if err != jwt.ErrExpiredToken {
			logger.Error(err)
		}
		jsonResponse(w, http.StatusUnauthorized, nil, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, jwtUser, "")
}

func changePassword(w http.ResponseWriter, r *http.Request) {
	jwtUser, err := jwt.Authenticate(r)
	if err != nil && err != jwt.ErrExpiredToken {
		if err != jwt.ErrExpiredToken {
			logger.Error(err)
		}
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

	err = user.ChangePassword(jwtUser.UserID, request.NewPassword)
	if err != nil {
		logger.Error(err)
		jsonResponse(w, http.StatusNotAcceptable, nil, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, jwtUser, "")
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	_, err := jwt.Authenticate(r)
	if err != nil && err != jwt.ErrExpiredToken {
		if err != jwt.ErrExpiredToken {
			logger.Error(err)
		}
		jsonResponse(w, http.StatusUnauthorized, nil, err.Error())
		return
	}

	users, err := user.GetAll()
	if err != nil {
		logger.Error(err)
		jsonResponse(w, http.StatusInternalServerError, nil, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, users, "")
}
