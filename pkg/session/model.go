package session

import (
	"time"

	"github.com/wexel-nath/wexel-auth/pkg/database"
)

type Session struct {
	SessionID string    `json:"session_id"`
	UserID    int64     `json:"user_id"`
	Created   time.Time `json:"session_created"`
	Expiry    time.Time `json:"session_expiry"`
}

func newSessionFromRow(row map[string]interface{}) (Session, error) {
	session := Session{}
	var ok bool

	if session.SessionID, ok = row[columnSessionID].(string); !ok {
		return session, database.RowError(row, columnSessionID, "string")
	}
	if session.UserID, ok = row[columnUserID].(int64); !ok {
		return session, database.RowError(row, columnUserID, "int64")
	}
	if session.Created, ok = row[columnCreated].(time.Time); !ok {
		return session, database.RowError(row, columnCreated, "time.Time")
	}
	if session.Expiry, ok = row[columnExpiry].(time.Time); !ok {
		return session, database.RowError(row, columnExpiry, "time.Time")
	}

	return session, nil
}
