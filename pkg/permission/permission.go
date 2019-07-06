package permission

import "github.com/wexel-nath/wexel-auth/pkg/logger"

func GetAllForUser(userID int64) (UserPermissions, error) {
	logger.Info("Getting all permissions for user[%d]", userID)

	rows, err := selectAllForUser(userID)
	if err != nil {
		return UserPermissions{}, err
	}

	return newUserPermissionsFromRows(rows)
}
