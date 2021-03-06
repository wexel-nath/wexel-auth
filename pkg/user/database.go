package user

import (
	"strings"

	"wexel-auth/pkg/database"
)

const (
	// Columns
	columnUserID    = "user_id"
	columnFirstName = "first_name"
	columnLastName  = "last_name"
	columnEmail     = "email"
	columnUsername  = "username"
	columnPassword  = "password"

	columnSessionCreated = "session_created"

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

	selectUserJoinSessionColumns = []string{
		columnUserID,
		columnFirstName,
		columnLastName,
		columnEmail,
		columnUsername,
		columnSessionCreated,
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

func selectByUsername(username string) (map[string]interface{}, error) {
	query := `
		SELECT
			` + strings.Join(selectUserColumns, ", ") + `
		FROM
			users
		WHERE
			` + columnUsername + ` = $1
	`

	db := database.GetConnection()
	row := db.QueryRow(query, username)
	return database.ScanRowToMap(row, selectUserColumns)
}

func updatePassword(userID int64, password string) (map[string]interface{}, error) {
	query := `
		UPDATE
			users
		SET
			` + columnPassword + ` = crypt($2, gen_salt('` + saltAlgorithm + `'))
		WHERE
			` + columnUserID + ` = $1
		RETURNING
			` + strings.Join(selectUserColumns, ", ")

	db := database.GetConnection()
	row := db.QueryRow(query, userID, password)
	return database.ScanRowToMap(row, selectUserColumns)
}

func selectAll() ([]map[string]interface{}, error) {
	query := `
		SELECT DISTINCT ON (` + columnUserID + `)
			` + strings.Join(selectUserJoinSessionColumns, ", ") + `
		FROM
			users
			LEFT JOIN session USING (` + columnUserID + `)
		ORDER BY
			` + columnUserID + `, ` + columnSessionCreated + ` DESC
	`

	db := database.GetConnection()
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	return database.ScanRowsToMap(rows, selectUserJoinSessionColumns)
}
