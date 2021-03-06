package session

import (
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"wexel-auth/pkg/config"
	"wexel-auth/pkg/database"
)

func TestInsert(t *testing.T) {
	config.Configure()

	now := time.Now()
	expiry := now.Add(30 * time.Minute)

	type args struct{
		sessionID string
		userID    int64
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
				sessionID: "test.session.token.1",
				userID:    1,
			},
			mock: database.Mock{
				ExpectRows: []database.MockRow{
					[]driver.Value{
						"test.session.token.1",
						int64(1),
						now,
						expiry,
					},
				},
			},
			want: want{
				row: map[string]interface{}{
					columnSessionID: "test.session.token.1",
					columnUserID:    int64(1),
					columnCreated:   now,
					columnExpiry:    expiry,
				},
				err: nil,
			},
		},
		"error": {
			args: args{
				sessionID: "test.session.token.2",
				userID:    2,
			},
			mock: database.Mock{
				ExpectErr: errors.New("connection error"),
			},
			want: want{
				row: nil,
				err: errors.New("connection error"),
			},
		},
	}

	expectedQuery := `
		INSERT INTO session \(
			session_id,
			user_id,
			session_created,
			session_expiry
		\)
		VALUES \(
			(.+)
		\)
		RETURNING
			session_id,
			user_id,
			session_created,
			session_expiry
	`

	for name, test := range tests {
		t.Run(name, func(st *testing.T) {
			dbMock := database.GetMockDB(st)
			query := dbMock.ExpectQuery(expectedQuery)

			if test.mock.ExpectErr != nil {
				query.WillReturnError(test.mock.ExpectErr)
			} else {
				mockRows := sqlmock.NewRows(sessionColumns)
				for _, row := range test.mock.ExpectRows {
					mockRows.AddRow(row...)
				}
				query.WillReturnRows(mockRows)
			}

			row, err := insert(test.sessionID, test.userID)

			assert.Equal(st, test.row, row)
			assert.Equal(st, test.err, err)
		})
	}
}

func TestUpdateSessionExpiry(t *testing.T) {
	now := time.Now()
	extension := 30 * time.Minute

	type args struct{
		sessionID string
		userID    int64
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
				sessionID: "test.session.token.1",
				userID:    1,
			},
			mock: database.Mock{
				ExpectRows: []database.MockRow{
					[]driver.Value{
						"test.session.token.1",
						int64(1),
						now,
						now.Add(extension),
					},
				},
			},
			want: want{
				row: map[string]interface{}{
					columnSessionID: "test.session.token.1",
					columnUserID:    int64(1),
					columnCreated:   now,
					columnExpiry:    now.Add(extension),
				},
				err: nil,
			},
		},
		"session not found": {
			args: args{
				sessionID: "test.session.token.2",
				userID:    2,
			},
			mock: database.Mock{
				ExpectErr: errors.New("session not found"),
			},
			want: want{
				row: nil,
				err: errors.New("session not found"),
			},
		},
	}

	expectedQuery := `
		UPDATE
			session
		SET
			session_expiry = \$1
		WHERE
			session_id = \$2
			AND user_id = \$3
			AND session_expiry > \$4
		RETURNING
			session_id,
			user_id,
			session_created,
			session_expiry
	`

	for name, test := range tests {
		t.Run(name, func(st *testing.T) {
			dbMock := database.GetMockDB(st)
			query := dbMock.ExpectQuery(expectedQuery)

			if test.mock.ExpectErr != nil {
				query.WillReturnError(test.mock.ExpectErr)
			} else {
				mockRows := sqlmock.NewRows(sessionColumns)
				for _, row := range test.mock.ExpectRows {
					mockRows.AddRow(row...)
				}
				query.WillReturnRows(mockRows)
			}

			row, err := updateSessionExpiry(test.sessionID, test.userID, extension)

			assert.Equal(st, test.row, row)
			assert.Equal(st, test.err, err)
		})
	}
}
