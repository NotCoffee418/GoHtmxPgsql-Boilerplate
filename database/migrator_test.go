package database

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_getInstalledMigrationVersion(t *testing.T) {
	db, mock := GetMockDB()

	// Mock for successful query
	rows := sqlmock.NewRows([]string{"version"}).AddRow(5)
	mock.ExpectQuery("SELECT version FROM migrations ORDER BY version DESC LIMIT 1").WillReturnRows(rows)

	// Channel to collect the result
	resultChan := make(chan int, 1)

	getInstalledMigrationVersion(db, resultChan)
	version := <-resultChan

	// Validate
	assert.Equal(t, 5, version)
}
