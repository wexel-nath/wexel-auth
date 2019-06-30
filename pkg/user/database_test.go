package user

import (
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/wexel-nath/wexel-auth/pkg/database"
)

func TestInsert(t *testing.T) {
	type args struct{
		firstName string
		lastName  string
		email     string
		username  string
		password  string
	}
	type want struct{
		row map[string]interface{}
		err error
	}
	tests := map[string]struct{
		args
		mock database.Mock
		want
	}{
		"success": {
			args: args{
				firstName: "Dev",
				lastName:  "User",
				email:     "dev.user@test.com",
				username:  "dev",
				password:  "4Me2Test",
			},
			mock: database.Mock{
				ExpectRows: []database.MockRow{
					[]driver.Value{
						int64(1),
						"Dev",
						"User",
						"dev.user@test.com",
						"dev",
					},
				},
			},
			want: want{
				row: map[string]interface{}{
					ColumnUserID:    int64(1),
					ColumnFirstName: "Dev",
					ColumnLastName:  "User",
					ColumnEmail:     "dev.user@test.com",
					ColumnUsername:  "dev",
				},
				err: nil,
			},
		},
		"bad username": {
			args: args{
				firstName: "Bad",
				lastName:  "User",
				email:     "bad.user@test.com",
				username:  "bad_username",
				password:  "4Me2Test",
			},
			mock: database.Mock{
				ExpectErr: errors.New("error: bad username"),
			},
			want: want{
				row: nil,
				err: errors.New("error: bad username"),
			},
		},
	}

	expectedQuery := `
		INSERT INTO users \(
			first_name,
			last_name,
			email,
			username,
			password
		\)
		VALUES \(
			(.+)
		\)
		RETURNING
			user_id,
			first_name,
			last_name,
			email,
			username
	`

	for name, test := range tests {
		t.Run(name, func(st *testing.T) {
			dbMock := database.GetMockDB(st)
			query := dbMock.ExpectQuery(expectedQuery)

			if test.mock.ExpectErr != nil {
				query.WillReturnError(test.mock.ExpectErr)
			} else {
				mockRows := sqlmock.NewRows(selectUserColumns)
				for _, row := range test.mock.ExpectRows {
					mockRows.AddRow(row...)
				}
				query.WillReturnRows(mockRows)
			}

			row, err := insert(test.firstName, test.lastName, test.email, test.username, test.password)

			assert.Equal(st, test.row, row)
			assert.Equal(st, test.err, err)
		})
	}
}

func TestSelectByCredentials(t *testing.T) {
	type args struct{
		username  string
		password  string
	}
	type want struct{
		row map[string]interface{}
		err error
	}
	tests := map[string]struct{
		args
		mock database.Mock
		want
	}{
		"success": {
			args: args{
				username:  "dev",
				password:  "4Me2Test",
			},
			mock: database.Mock{
				ExpectRows: []database.MockRow{
					[]driver.Value{
						int64(1),
						"Dev",
						"User",
						"dev.user@test.com",
						"dev",
					},
				},
			},
			want: want{
				row: map[string]interface{}{
					ColumnUserID:    int64(1),
					ColumnFirstName: "Dev",
					ColumnLastName:  "User",
					ColumnEmail:     "dev.user@test.com",
					ColumnUsername:  "dev",
				},
				err: nil,
			},
		},
		"bad username": {
			args: args{
				username:  "bad_username",
				password:  "4Me2Test",
			},
			mock: database.Mock{
				ExpectErr: errors.New("error: bad username"),
			},
			want: want{
				row: nil,
				err: errors.New("error: bad username"),
			},
		},
	}

	expectedQuery := `
		SELECT
			user_id,
			first_name,
			last_name,
			email,
			username
		FROM
			users
		WHERE
			username = (.+)
			AND password = crypt\((.+), password\)
	`

	for name, test := range tests {
		t.Run(name, func(st *testing.T) {
			dbMock := database.GetMockDB(st)
			query := dbMock.ExpectQuery(expectedQuery)

			if test.mock.ExpectErr != nil {
				query.WillReturnError(test.mock.ExpectErr)
			} else {
				mockRows := sqlmock.NewRows(selectUserColumns)
				for _, row := range test.mock.ExpectRows {
					mockRows.AddRow(row...)
				}
				query.WillReturnRows(mockRows)
			}

			row, err := selectByCredentials(test.username, test.password)

			assert.Equal(st, test.row, row)
			assert.Equal(st, test.err, err)
		})
	}
}
