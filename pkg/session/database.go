package session

import (
	"strings"
	"time"

	"github.com/wexel-nath/wexel-auth/pkg/config"
	"github.com/wexel-nath/wexel-auth/pkg/database"
)

const (
	// Columns
	columnSessionID = "session_id"
	columnUserID    = "user_id"
	columnTimestamp = "timestamp"
	columnExpiry    = "expiry"
)

var (
	sessionColumns = []string{
		columnSessionID,
		columnUserID,
		columnTimestamp,
		columnExpiry,
	}
)

func insert(sessionID string, userID int64) (map[string]interface{}, error) {
	query := `
		INSERT INTO session (
			` + strings.Join(sessionColumns, ", ") + `
		)
		VALUES (
			$1,
			$2,
			$3,
			$4
		)
		RETURNING
			` + strings.Join(sessionColumns, ", ")

	timestamp := time.Now().Unix()
	expiry := timestamp + config.GetSessionExpiry()

	db := database.GetConnection()
	row := db.QueryRow(query, sessionID, userID, timestamp, expiry)
	return database.ScanRowToMap(row, sessionColumns)
}

func selectActiveSession(sessionID string, userID int64) (map[string]interface{}, error) {
	query := `
		SELECT
			` + strings.Join(sessionColumns, ", ") + `
		FROM
			session
		WHERE
			` + columnSessionID + ` = $1
			AND ` + columnUserID + ` = $2
			AND ` + columnExpiry + ` > $3
	`

	now := time.Now().Unix()

	db := database.GetConnection()
	row := db.QueryRow(query, sessionID, userID, now)
	return database.ScanRowToMap(row, sessionColumns)
}

func updateSessionExpiry(sessionID string, userID int64, extension int64) (map[string]interface{}, error) {
	query := `
		UPDATE
			session
		SET
			` + columnExpiry + ` = $1
		WHERE
			` + columnSessionID + ` = $2
			AND ` + columnUserID + ` = $3
			AND ` + columnExpiry + ` > $4
		RETURNING
			` + strings.Join(sessionColumns, ", ")

	now := time.Now().Unix()
	newExpiry := now + extension

	db := database.GetConnection()
	row := db.QueryRow(query, newExpiry, sessionID, userID, now)
	return database.ScanRowToMap(row, sessionColumns)
}
