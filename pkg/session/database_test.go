package session

import (
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/wexel-nath/wexel-auth/pkg/database"
)

func TestInsert(t *testing.T) {
	now := time.Now().Unix()

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
						now + 300,
					},
				},
			},
			want: want{
				row: map[string]interface{}{
					columnSessionID: "test.session.token.1",
					columnUserID:    int64(1),
					columnTimestamp: now,
					columnExpiry:    now + 300,
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
			timestamp,
			expiry
		\)
		VALUES \(
			(.+)
		\)
		RETURNING
			session_id,
			user_id,
			timestamp,
			expiry
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
	now := time.Now().Unix()

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
						now + 300,
					},
				},
			},
			want: want{
				row: map[string]interface{}{
					columnSessionID: "test.session.token.1",
					columnUserID:    int64(1),
					columnTimestamp: now,
					columnExpiry:    now + 300,
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
			expiry = \$1
		WHERE
			session_id = \$2
			AND user_id = \$3
			AND expiry > \$4
		RETURNING
			session_id,
			user_id,
			timestamp,
			expiry
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

			row, err := updateSessionExpiry(test.sessionID, test.userID)

			assert.Equal(st, test.row, row)
			assert.Equal(st, test.err, err)
		})
	}
}
