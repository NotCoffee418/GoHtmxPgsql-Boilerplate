package server

import (
	"fmt"
	"sync"

	"github.com/NotCoffee418/GoWebsite-Boilerplate/internal/access"
	"github.com/NotCoffee418/GoWebsite-Boilerplate/internal/config"
	log "github.com/sirupsen/logrus"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/NotCoffee418/GoWebsite-Boilerplate/internal/utils"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	once sync.Once
	db   *sqlx.DB
)

// GetDB returns a singleton DB instance
func GetDB() *sqlx.DB {
	once.Do(func() {
		db = initDb()
	})
	return db
}

// GetMockDB returns a mock DB for testing
func GetMockDB() (*sqlx.DB, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return sqlx.NewDb(mockDB, "sqlmock"), mock
}

func initDb() *sqlx.DB {
	// Get connection details from environment
	pgHost := utils.GetEnv("PG_HOST")
	pgPort := utils.GetEnv("PG_PORT")
	pgUser := utils.GetEnv("PG_USER")
	pgPass := utils.GetEnv("PG_PASS")
	pgDbName := utils.GetEnv("PG_DATABASE")
	pgSslMode := utils.GetEnv("PG_SSL_MODE")

	// Validate SSL mode
	validSslModes := []string{
		"disable", "allow", "prefer", "require", "verify-ca", "verify-full"}
	if !utils.SliceContains(validSslModes, pgSslMode) {
		log.Fatalf("Invalid PG_SSL_MODE: %s", pgSslMode)
	}

	// Set up connection
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		pgHost, pgPort, pgUser, pgPass, pgDbName, pgSslMode)
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to DB: %v", err)
	}

	// Set up pooling rules
	db.SetMaxOpenConns(config.DbMaxOpenConns)
	db.SetMaxIdleConns(config.DbMaxIdleConns)
	db.SetConnMaxLifetime(config.DbConnMaxLifetime)

	access.DB = db
	return db
}
