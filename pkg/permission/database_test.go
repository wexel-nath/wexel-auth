package permission

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/wexel-nath/wexel-auth/pkg/database"
)

func TestSelectAllForUser(t *testing.T) {
	type args struct{
		userID    int64
	}
	type want struct{
		rows []map[string]interface{}
		err  error
	}
	tests := map[string]struct{
		args
		mock database.Mock
		want
	}{
		"now rows found": {
			args: args{
				userID: 2,
			},
			mock: database.Mock{
				ExpectErr: errors.New("session not found"),
			},
			want: want{
				rows: nil,
				err: errors.New("session not found"),
			},
		},
	}

	expectedQuery := `
		SELECT
			service_name,
			ARRAY_AGG\(permission_name\) permissions
		FROM
			service
			JOIN permission USING \(service_id\)
			JOIN user_permission USING \(permission_id\)
		WHERE
			user_id = \$1
		GROUP BY
			service_name
	`

	for name, test := range tests {
		t.Run(name, func(st *testing.T) {
			dbMock := database.GetMockDB(st)
			query := dbMock.ExpectQuery(expectedQuery)

			if test.mock.ExpectErr != nil {
				query.WillReturnError(test.mock.ExpectErr)
			} else {
				mockRows := sqlmock.NewRows(selectAllColumns)
				for _, row := range test.mock.ExpectRows {
					mockRows.AddRow(row...)
				}
				query.WillReturnRows(mockRows)
			}

			rows, err := selectAllForUser(test.userID)

			assert.Equal(st, test.rows, rows)
			assert.Equal(st, test.err, err)
		})
	}
}
