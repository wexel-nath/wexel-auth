package session

import (
	"fmt"
	"time"
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
		return session, fmt.Errorf("row[%v] does not contain field[%s] type[string]", row, columnSessionID)
	}
	if session.UserID, ok = row[columnUserID].(int64); !ok {
		return session, fmt.Errorf("row[%v] does not contain field[%s] type int64", row, columnUserID)
	}
	if session.Created, ok = row[columnCreated].(time.Time); !ok {
		return session, fmt.Errorf("row[%v] does not contain field[%s] type time.Time", row, columnCreated)
	}
	if session.Expiry, ok = row[columnExpiry].(time.Time); !ok {
		return session, fmt.Errorf("row[%v] does not contain field[%s] type time.Time", row, columnExpiry)
	}

	return session, nil
}
