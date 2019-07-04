package session

import "fmt"

type Session struct {
	SessionID string `json:"session_id"`
	UserID    int64  `json:"user_id"`
	Timestamp int64  `json:"timestamp"`
	Expiry    int64  `json:"expiry"`
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
	if session.Timestamp, ok = row[columnTimestamp].(int64); !ok {
		return session, fmt.Errorf("row[%v] does not contain field[%s] type int64", row, columnTimestamp)
	}
	if session.Expiry, ok = row[columnExpiry].(int64); !ok {
		return session, fmt.Errorf("row[%v] does not contain field[%s] type int64", row, columnExpiry)
	}

	return session, nil
}
