package user

import (
	"strings"

	"github.com/wexel-nath/wexel-auth/pkg/database"
)

const (
	// Columns
	columnUserID    = "user_id"
	columnFirstName = "first_name"
	columnLastName  = "last_name"
	columnEmail     = "email"
	columnUsername  = "username"
	columnPassword  = "password"

	// Crypto Salt Algorithm
	saltAlgorithm = "bf"
)

var (
	insertUserColumns = []string{
		columnFirstName,
		columnLastName,
		columnEmail,
		columnUsername,
		columnPassword,
	}

	selectUserColumns = []string{
		columnUserID,
		columnFirstName,
		columnLastName,
		columnEmail,
		columnUsername,
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
			` + columnUsername + ` = $1
			AND ` + columnPassword + ` = crypt($2, ` + columnPassword + `)
	`

	db := database.GetConnection()
	row := db.QueryRow(query, username, password)
	return database.ScanRowToMap(row, selectUserColumns)
}
