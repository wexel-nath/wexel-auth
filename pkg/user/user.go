package user

import (
	"database/sql"
	"errors"
	"fmt"
	"unicode"

	"github.com/wexel-nath/wexel-auth/pkg/logger"
)

const (
	weakPasswordMessage = "Your new password must be at least 8 characters, and include at least one number, lower case, and upper case letter."
)

var (
	ErrInvalidDetails = errors.New("Invalid username or password")
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
		return User{}, ErrInvalidDetails
	}
	if err != nil {
		logger.Error(err)
		return User{}, ErrInvalidDetails
	}

	return newUserFromRow(row)
}

func ChangePassword(userID int64, password string) error {
	logger.Info("Changing password for user[%d]", userID)

	if !isValid(password) {
		return fmt.Errorf(weakPasswordMessage)
	}

	_, err := updatePassword(userID, password)
	return err
}

func isValid(password string) bool {
	var hasNumber, hasLower, hasUpper bool
	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsUpper(c):
			hasUpper = true
		}
	}
	return hasNumber && hasLower && hasUpper && len(password) >= 8
}
