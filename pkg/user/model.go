package user

import (
	"github.com/wexel-nath/wexel-auth/pkg/database"
)

// User represents a row of the user table
type User struct {
	UserID    int64  `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Username  string `json:"username"`
}

func newUser(row map[string]interface{}) (User, error) {
	user := User{}
	var ok bool

	if user.UserID, ok = row[columnUserID].(int64); !ok {
		return user, database.RowError(row, columnUserID, "int64")
	}
	if user.FirstName, ok = row[columnFirstName].(string); !ok {
		return user, database.RowError(row, columnFirstName, "string")
	}
	if user.LastName, ok = row[columnLastName].(string); !ok {
		return user, database.RowError(row, columnLastName, "string")
	}
	if user.Email, ok = row[columnEmail].(string); !ok {
		return user, database.RowError(row, columnEmail, "string")
	}
	if user.Username, ok = row[columnUsername].(string); !ok {
		return user, database.RowError(row, columnUsername, "string")
	}

	return user, nil
}

func buildFromRow(row map[string]interface{}, err error) (User, error) {
	if err != nil {
		return User{}, err
	}
	return newUser(row)
}

func buildFromRows(rows []map[string]interface{}, err error) ([]User, error) {
	users := make([]User, 0)
	if err != nil {
		return users, err
	}

	for _, row := range rows {
		u, err := newUser(row)
		if err != nil {
			return users, err
		}
		users = append(users, u)
	}

	return users, nil
}
