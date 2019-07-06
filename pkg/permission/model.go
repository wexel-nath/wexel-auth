package permission

import (
	"fmt"
)

type UserPermissions map[string][]string

func newUserPermissionsFromRows(rows []map[string]interface{}) (UserPermissions, error) {
	userPermissions := UserPermissions{}

	for _, row := range rows {
		serviceName, ok := row[columnServiceName].(string)
		if !ok {
			return userPermissions, fmt.Errorf("row[%v] does not contain field[%s] type[string]", row, columnServiceName)
		}
		permissions, ok := row[columnPermissions].([]string)
		if !ok {
			return userPermissions, fmt.Errorf("row[%v] does not contain field[%s] type[[]string]", row, columnPermissions)
		}

		userPermissions[serviceName] = permissions
	}

	return userPermissions, nil
}
