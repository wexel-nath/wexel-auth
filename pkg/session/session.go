package session

import "github.com/wexel-nath/wexel-auth/pkg/logger"

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
	if err != nil {
		return Session{}, err
	}

	return newSessionFromRow(row)
}

func ExtendCurrentSession(sessionID string, userID int64) (Session, error) {
	logger.Info("Updating current session[%s] for user[%d]", sessionID, userID)

	row, err := updateSessionExpiry(sessionID, userID)
	if err != nil {
		return Session{}, err
	}

	return newSessionFromRow(row)
}