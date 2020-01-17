package user

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/wexel-nath/wexel-auth/pkg/logger"
	"github.com/wexel-nath/wexel-auth/pkg/util"
)

const (
	defaultPassword = "4Me2Change"
	weakPasswordMessage = "Your new password must be at least 8 characters, and include at least one number, lower case, and upper case letter."
)

var (
	ErrInvalidDetails = errors.New("Invalid username or password")
)

func Create(
	firstName string,
	lastName string,
	email string,
) (User, error) {
	logger.Info("Creating user[%s %s]", firstName, lastName)

	username, err := generateUsername(firstName, lastName)
	if err != nil {
		return User{}, fmt.Errorf("creating user[%s] failed: %v", username, err)
	}

	row, err := insert(firstName, lastName, email, username, defaultPassword)
	if err != nil {
		return User{}, fmt.Errorf("creating user[%s] failed: %v", username, err)
	}

	return newUser(row)
}

// generates a username of the form 'flast', if taken then 'filast'
func generateUsername(firstName string, lastName string) (string, error) {
	first := strings.ToLower(util.StripWhitespace(firstName))
	last := strings.ToLower(util.StripWhitespace(lastName))
	firstPart := ""
	for _, c := range first {
		firstPart += string(c)
		username := firstPart + last
		taken, err := isUsernameTaken(firstPart + last)
		if err != nil {
			return "", err
		}
		if !taken {
			return username, nil
		}
	}

	return "", fmt.Errorf("could not find an available username for [%s %s]", firstName, lastName)
}

func isUsernameTaken(username string) (bool, error) {
	_, err := selectByUsername(username)
	if err == sql.ErrNoRows {
		return false, nil
	}

	return true, err
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

	return newUser(row)
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

func GetAll() ([]User, error) {
	rows, err := selectAll()
	return buildFromRows(rows, err)
}
