package user

import (
	"fmt"
)

// User represents a row of the user table
type User struct {
	UserID    int64  `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Username  string `json:"username"`
}

func newUserFromRow(row map[string]interface{}) (User, error) {
	user := User{}
	var ok bool

	if user.UserID, ok = row[columnUserID].(int64); !ok {
		return user, fmt.Errorf("row[%v] does not contain field[%s] type[int64]", row, columnUserID)
	}
	if user.FirstName, ok = row[columnFirstName].(string); !ok {
		return user, fmt.Errorf("row[%v] does not contain field[%s] type[string]", row, columnFirstName)
	}
	if user.LastName, ok = row[columnLastName].(string); !ok {
		return user, fmt.Errorf("row[%v] does not contain field[%s] type[string]", row, columnLastName)
	}
	if user.Email, ok = row[columnEmail].(string); !ok {
		return user, fmt.Errorf("row[%v] does not contain field[%s] type[string]", row, columnEmail)
	}
	if user.Username, ok = row[columnUsername].(string); !ok {
		return user, fmt.Errorf("row[%v] does not contain field[%s] type[string]", row, columnUsername)
	}

	return user, nil
}
