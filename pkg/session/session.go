package session

import (
	"database/sql"
	"errors"

	"wexel-auth/pkg/config"
	"wexel-auth/pkg/logger"
)

var (
	ErrSessionExpired = errors.New("session has expired")
)

func Create(userID int64) (Session, error) {
	logger.Info("Creating session for user[%d]", userID)

	sessionID, err := generateUniqueID(userID)
	if err != nil {
		return Session{}, err
	}

	row, err := insert(sessionID, userID)
	if err != nil {
		return Session{}, err
	}

	return newSessionFromRow(row)
}

func GetCurrentSession(sessionID string, userID int64) (Session, error) {
	logger.Info("Getting current session[%s] for user[%d]", sessionID, userID)

	row, err := selectActiveSession(sessionID, userID)
	if err == sql.ErrNoRows {
		return Session{}, ErrSessionExpired
	}
	if err != nil {
		return Session{}, err
	}

	return newSessionFromRow(row)
}

func ExtendCurrentSession(sessionID string, userID int64) (Session, error) {
	logger.Info("Updating current session[%s] for user[%d]", sessionID, userID)

	row, err := updateSessionExpiry(sessionID, userID, config.GetSessionExpiry())
	if err == sql.ErrNoRows {
		return Session{}, ErrSessionExpired
	}
	if err != nil {
		return Session{}, err
	}

	return newSessionFromRow(row)
}

func EndCurrentSession(sessionID string, userID int64) (Session, error) {
	logger.Info("Ending current session[%s] for user[%d]", sessionID, userID)

	row, err := updateSessionExpiry(sessionID, userID, 0)
	if err == sql.ErrNoRows {
		return Session{}, ErrSessionExpired
	}
	if err != nil {
		return Session{}, err
	}

	return newSessionFromRow(row)
}
