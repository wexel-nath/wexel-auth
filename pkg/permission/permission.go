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

func AddUserPermission(userID int64, permission string) error {
	logger.Info("Adding permission[%s] for user [%d]", permission, userID)

	return insertUserPermissionByName(userID, permission)
}
