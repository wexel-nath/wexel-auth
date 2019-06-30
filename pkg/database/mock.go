package database

import (
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

type MockRow []driver.Value

type Mock struct {
	ExpectRows []MockRow
	ExpectErr  error
}

func GetMockDB(t *testing.T) sqlmock.Sqlmock {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	connection = db
	return mock
}
