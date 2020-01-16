package database

import "fmt"

func RowError(row map[string]interface{}, field string, castType string) error {
	return fmt.Errorf(
		"row[%v] does not contain field[%s] type[%s]",
		row,
		field,
		castType,
	)
}
