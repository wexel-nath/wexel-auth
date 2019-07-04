package user

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/wexel-nath/wexel-auth/pkg/logger"
)

func Create(
	firstName string,
	lastName string,
	email string,
	username string,
	password string,
) (User, error) {
	logger.Info("Creating user[%s]", username)

	row, err := insert(firstName, lastName, email, username, password)
	if err != nil {
		return User{}, fmt.Errorf("creating user[%s] failed: %v", username, err)
	}

	return newUserFromRow(row)
}

func Authenticate(username string, password string) (User, error) {
	logger.Info("Authenticating user[%s]", username)

	row, err := selectByCredentials(username, password)
	if err == sql.ErrNoRows {
		return User{}, errors.New("Invalid username or password")
	}
	if err != nil {
		return User{}, fmt.Errorf("authenticating user[%s] failed: %v", username, err)
	}

	return newUserFromRow(row)
}
