package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/wexel-nath/wexel-auth/pkg/user"
)

func HandleCreateUser(r *http.Request) (interface{}, int, error) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return nil, http.StatusBadRequest, err
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
		return nil, http.StatusBadRequest, err
	}

	userModel, err := user.Create(
		request.FirstName,
		request.LastName,
		request.Email,
		request.Username,
		request.Password,
	)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	return userModel, http.StatusCreated, nil
}
