package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUserFromRow(t *testing.T) {
	type want struct{
		user User
		wantErr  bool
	}
	tests := map[string]struct{
		row  map[string]interface{}
		want
	}{
		"success": {
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
				wantErr: false,
			},
		},
		"missing id": {
			row:  map[string]interface{}{
				ColumnFirstName: "Dev",
				ColumnLastName:  "User",
				ColumnEmail:     "dev.user@test.com",
				ColumnUsername:  "dev",
			},
			want: want{
				wantErr: true,
			},
		},
		"missing first name": {
			row:  map[string]interface{}{
				ColumnUserID:    int64(1),
				ColumnLastName:  "User",
				ColumnEmail:     "dev.user@test.com",
				ColumnUsername:  "dev",
			},
			want: want{
				wantErr: true,
			},
		},
		"missing last name": {
			row:  map[string]interface{}{
				ColumnUserID:    int64(1),
				ColumnFirstName: "Dev",
				ColumnEmail:     "dev.user@test.com",
				ColumnUsername:  "dev",
			},
			want: want{
				wantErr: true,
			},
		},
		"missing email": {
			row:  map[string]interface{}{
				ColumnUserID:    int64(1),
				ColumnFirstName: "Dev",
				ColumnLastName:  "User",
				ColumnUsername:  "dev",
			},
			want: want{
				wantErr: true,
			},
		},
		"missing username": {
			row:  map[string]interface{}{
				ColumnUserID:    int64(1),
				ColumnFirstName: "Dev",
				ColumnLastName:  "User",
				ColumnEmail:     "dev.user@test.com",
			},
			want: want{
				wantErr: true,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(st *testing.T) {
			user, err := newUserFromRow(test.row)

			assert.Equal(t, test.wantErr, err != nil)
			if !test.wantErr {
				assert.Equal(t, test.user, user)
			}
		})
	}
}
