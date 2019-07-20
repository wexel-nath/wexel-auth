package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/wexel-nath/wexel-auth/pkg/logger"
	"github.com/wexel-nath/wexel-auth/pkg/user"
)

func createUser(r *http.Request, _ interface{}) (interface{}, interface{}, int) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		logger.Error(err)
		return nil, err.Error(), http.StatusBadRequest
	}

	var request struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Username  string `json:"username"`
		Password  string `json:"password"`
	}
	err = json.Unmarshal(body, &request)
	if err != nil {
		logger.Error(err)
		return nil, err.Error(), http.StatusBadRequest
	}

	userModel, err := user.Create(
		request.FirstName,
		request.LastName,
		request.Email,
		request.Username,
		request.Password,
	)
	if err != nil {
		logger.Error(err)
		return nil, err.Error(), http.StatusBadRequest
	}
	return userModel, nil, http.StatusCreated
}

func getUser(_ *http.Request, user interface{}) (interface{}, interface{}, int) {
	return user, nil, http.StatusOK
}
