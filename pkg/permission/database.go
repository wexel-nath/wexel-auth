package permission

import (
	"github.com/lib/pq"
	"github.com/wexel-nath/wexel-auth/pkg/database"
)

const (
	// Columns
	columnPermissionID = "permission_id"
	columnPermissionName = "permission_name"
	columnPermissions = "permissions"
	columnServiceID = "service_id"
	columnServiceName = "service_name"
	columnUserID = "user_id"
)

var (
	selectAllColumns = []string{
		columnServiceName,
		columnPermissions,
	}
)

func insertUserPermissions(userID int64, permissionIDs []int64) error {
	query := `
		INSERT INTO user_permission (
			` + columnUserID + `,
			` + columnPermissionID + `
		)
			SELECT
				$1,
				UNNEST($2)
	`

	db := database.GetConnection()
	_, err := db.Exec(query, userID, pq.Array(permissionIDs))
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
