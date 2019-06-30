package user

import (
	"strings"

	"github.com/wexel-nath/wexel-auth/pkg/database"
)

const (
	// Columns
	ColumnUserID    = "user_id"
	ColumnFirstName = "first_name"
	ColumnLastName  = "last_name"
	ColumnEmail     = "email"
	ColumnUsername  = "username"
	ColumnPassword  = "password"

	// Crypto Salt Algorithm
	saltAlgorithm = "bf"
)

var (
	insertUserColumns = []string{
		ColumnFirstName,
		ColumnLastName,
		ColumnEmail,
		ColumnUsername,
		ColumnPassword,
	}

	selectUserColumns = []string{
		ColumnUserID,
		ColumnFirstName,
		ColumnLastName,
		ColumnEmail,
		ColumnUsername,
	}
)

func insert(
	firstName string,
	lastName string,
	email string,
	username string,
	password string,
) (map[string]interface{}, error) {
	query := `
		INSERT INTO users (
			` + strings.Join(insertUserColumns, ", ") + `
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			crypt($5, gen_salt('` + saltAlgorithm + `'))
		)
		RETURNING
			` + strings.Join(selectUserColumns, ", ")

	db := database.GetConnection()
	row := db.QueryRow(query, firstName, lastName, email, username, password)
	return database.ScanRowToMap(row, selectUserColumns)
}

func selectByCredentials(username string, password string) (map[string]interface{}, error) {
	query := `
		SELECT
			` + strings.Join(selectUserColumns, ", ") + `
		FROM
			users
		WHERE
			` + ColumnUsername + ` = $1
			AND ` + ColumnPassword + ` = crypt($2, ` + ColumnPassword + `)
	`

	db := database.GetConnection()
	row := db.QueryRow(query, username, password)
	return database.ScanRowToMap(row, selectUserColumns)
}
