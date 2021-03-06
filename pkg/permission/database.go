package permission

import (
	"strings"

	"github.com/lib/pq"
	"wexel-auth/pkg/database"
)

const (
	// Columns
	columnPermissionID          = "permission_id"
	columnPermissionName        = "permission_name"
	columnPermissionDescription = "permission_description"
	columnPermissions           = "permissions"
	columnServiceID             = "service_id"
	columnServiceName           = "service_name"
	columnUserID                = "user_id"
)

var (
	selectPermissionsColumns = []string{
		columnPermissionID,
		columnPermissionName,
		columnPermissionDescription,
	}

	selectAllColumns = []string{
		columnServiceName,
		columnPermissions,
	}

	insertUserPermissionColumns = []string{
		columnUserID,
		columnPermissionID,
	}
)

func selectAllForService(serviceName string) ([]map[string]interface{}, error) {
	query := `
		SELECT
			` + strings.Join(selectPermissionsColumns, ", ") + `
		FROM
			permission
			JOIN service USING (` + columnServiceID + `)
		WHERE
			` + columnServiceName + ` = $1
			OR ` + columnServiceName + ` = 'all'
	`

	db := database.GetConnection()
	rows, err := db.Query(query, serviceName)
	if err != nil {
		return nil, err
	}
	return database.ScanRowsToMap(rows, selectPermissionsColumns)
}

func insertUserPermissions(userID int64, permissionIDs []int64) error {
	query := `
		INSERT INTO user_permission (
			` + strings.Join(insertUserPermissionColumns, ", ") + `
		)
		SELECT
			$1,
			UNNEST($2::INTEGER[])
	`

	db := database.GetConnection()
	_, err := db.Exec(query, userID, pq.Array(permissionIDs))
	return err
}

func insertUserPermissionByName(userID int64, permission string) error {
	query := `
		INSERT INTO user_permission (
			` + strings.Join(insertUserPermissionColumns, ", ") + `
		)
		SELECT
			$1,
			` + columnPermissionID + `
		FROM
			permission
		WHERE
			` + columnPermissionName + ` = $2
	`

	db := database.GetConnection()
	_, err := db.Exec(query, userID, permission)
	return err
}

func selectAllForUser(userID int64) ([]map[string]interface{}, error) {
	query := `
		SELECT
			` + columnServiceName + `,
			ARRAY_AGG(` + columnPermissionName + `) ` + columnPermissions + `
		FROM
			service
			JOIN permission USING (` + columnServiceID + `)
			JOIN user_permission USING (` + columnPermissionID + `)
		WHERE
			` + columnUserID + ` = $1
		GROUP BY
			` + columnServiceName

	db := database.GetConnection()
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0)

	// manually scan through rows to use pq.Array
	for rows.Next() {
		var serviceName string
		var permissions []string
		err = rows.Scan(&serviceName, pq.Array(&permissions))
		if err != nil {
			return nil, err
		}
		result = append(
			result,
			map[string]interface{}{
				columnServiceName: serviceName,
				columnPermissions: permissions,
			},
		)
	}

	return result, err
}
