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

	if user.UserID, ok = row[ColumnUserID].(int64); !ok {
		return user, fmt.Errorf("row[%v] does not contain field[%s] type[int64]", row, ColumnUserID)
	}
	if user.FirstName, ok = row[ColumnFirstName].(string); !ok {
		return user, fmt.Errorf("row[%v] does not contain field[%s] type[string]", row, ColumnFirstName)
	}
	if user.LastName, ok = row[ColumnLastName].(string); !ok {
		return user, fmt.Errorf("row[%v] does not contain field[%s] type[string]", row, ColumnLastName)
	}
	if user.Email, ok = row[ColumnEmail].(string); !ok {
		return user, fmt.Errorf("row[%v] does not contain field[%s] type[string]", row, ColumnEmail)
	}
	if user.Username, ok = row[ColumnUsername].(string); !ok {
		return user, fmt.Errorf("row[%v] does not contain field[%s] type[string]", row, ColumnUsername)
	}

	return user, nil
}
