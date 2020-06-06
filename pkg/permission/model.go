package permission

import (
	"wexel-auth/pkg/database"
)

type UserPermissions map[string][]string

func newUserPermissions(rows []map[string]interface{}) (UserPermissions, error) {
	userPermissions := UserPermissions{}

	for _, row := range rows {
		serviceName, ok := row[columnServiceName].(string)
		if !ok {
			return userPermissions, database.RowError(row, columnServiceName, "string")
		}
		permissions, ok := row[columnPermissions].([]string)
		if !ok {
			return userPermissions, database.RowError(row, columnPermissions, "[]string")
		}

		userPermissions[serviceName] = permissions
	}

	return userPermissions, nil
}

type Permission struct {
	ID          int64  `json:"permission_id"`
	Name        string `json:"permission_name"`
	Description string `json:"permission_description"`
}

func newPermission(row map[string]interface{}) (Permission, error) {
	permission := Permission{}
	var ok bool

	if permission.ID, ok = row[columnPermissionID].(int64); !ok {
		return permission, database.RowError(row, columnPermissionID, "int64")
	}
	if permission.Name, ok = row[columnPermissionName].(string); !ok {
		return permission, database.RowError(row, columnPermissionName, "string")
	}
	if permission.Description, ok = row[columnPermissionDescription].(string); !ok {
		return permission, database.RowError(row, columnPermissionDescription, "string")
	}

	return permission, nil
}

func buildFromRow(row map[string]interface{}, err error) (Permission, error) {
	if err != nil {
		return Permission{}, err
	}
	return newPermission(row)
}

func buildFromRows(rows []map[string]interface{}, err error) ([]Permission, error) {
	permissions := make([]Permission, 0)
	if err != nil {
		return permissions, err
	}

	for _, row := range rows {
		p, err := newPermission(row)
		if err != nil {
			return permissions, err
		}
		permissions = append(permissions, p)
	}

	return permissions, nil
}
