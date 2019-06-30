package user

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUserFromRow(t *testing.T) {
	type want struct{
		user User
		err  bool
	}
	tests := []struct{
		name string
		row  map[string]interface{}
		want
	}{
		{
			name: "success",
			row:  map[string]interface{}{
				ColumnUserID:    int64(1),
				ColumnFirstName: "Dev",
				ColumnLastName:  "User",
				ColumnEmail:     "dev.user@test.com",
				ColumnUsername:  "dev",
			},
			want: want{
				user: User{
					UserID:    1,
					FirstName: "Dev",
					LastName:  "User",
					Email:     "dev.user@test.com",
					Username:  "dev",
				},
				err: false,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(st *testing.T) {
			user, err := newUserFromRow(test.row)
			assert.Equal(t, test.user, user)
			assert.Equal(t, test.err, err != nil)
		})
	}
}
