package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUserFromRow(t *testing.T) {
	type want struct{
		user    User
		wantErr bool
	}
	tests := map[string]struct{
		row  map[string]interface{}
		want
	}{
		"success": {
			row:  map[string]interface{}{
				columnUserID:    int64(1),
				columnFirstName: "Dev",
				columnLastName:  "User",
				columnEmail:     "dev.user@test.com",
				columnUsername:  "dev",
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
				columnFirstName: "Dev",
				columnLastName:  "User",
				columnEmail:     "dev.user@test.com",
				columnUsername:  "dev",
			},
			want: want{
				wantErr: true,
			},
		},
		"missing first name": {
			row:  map[string]interface{}{
				columnUserID:    int64(1),
				columnLastName:  "User",
				columnEmail:     "dev.user@test.com",
				columnUsername:  "dev",
			},
			want: want{
				wantErr: true,
			},
		},
		"missing last name": {
			row:  map[string]interface{}{
				columnUserID:    int64(1),
				columnFirstName: "Dev",
				columnEmail:     "dev.user@test.com",
				columnUsername:  "dev",
			},
			want: want{
				wantErr: true,
			},
		},
		"missing email": {
			row:  map[string]interface{}{
				columnUserID:    int64(1),
				columnFirstName: "Dev",
				columnLastName:  "User",
				columnUsername:  "dev",
			},
			want: want{
				wantErr: true,
			},
		},
		"missing username": {
			row:  map[string]interface{}{
				columnUserID:    int64(1),
				columnFirstName: "Dev",
				columnLastName:  "User",
				columnEmail:     "dev.user@test.com",
			},
			want: want{
				wantErr: true,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(st *testing.T) {
			user, err := newUser(test.row)

			assert.Equal(t, test.wantErr, err != nil)
			if !test.wantErr {
				assert.Equal(t, test.user, user)
			}
		})
	}
}
