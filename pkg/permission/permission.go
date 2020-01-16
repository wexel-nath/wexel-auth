package permission

import (
	"github.com/wexel-nath/wexel-auth/pkg/logger"
)

func GetAllForUser(userID int64) (UserPermissions, error) {
	logger.Info("Getting all permissions for user[%d]", userID)

	rows, err := selectAllForUser(userID)
	if err != nil {
		return UserPermissions{}, err
	}

	return newUserPermissions(rows)
}

func AddUserPermissions(userID int64, permissions []int64) error {
	logger.Info("Adding permissions %v for user [%d]", permissions, userID)

	return insertUserPermissions(userID, permissions)
}

func GetAllForService(serviceName string) ([]Permission, error) {
	logger.Info("Getting permissions for service[%s]", serviceName)

	return buildFromRows(selectAllForService(serviceName))
}
