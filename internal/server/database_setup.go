package server

import (
	"fmt"
	"log"
	"sync"

	"github.com/NotCoffee418/GoWebsite-Boilerplate/config"
	"github.com/NotCoffee418/GoWebsite-Boilerplate/internal/common"
	"github.com/NotCoffee418/GoWebsite-Boilerplate/internal/utils"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	once sync.Once
	db   *sqlx.DB
)

func GetDBConn() *sqlx.DB {
	once.Do(func() {
		db = initDb()
	})
	return db
}

func initDb() *sqlx.DB {
	// Get connection details from environment
	pgHost := common.GetEnv("PG_HOST")
	pgPort := common.GetEnv("PG_PORT")
	pgUser := common.GetEnv("PG_USER")
	pgPass := common.GetEnv("PG_PASS")
	pgDbName := common.GetEnv("PG_DATABASE")
	pgSslMode := common.GetEnv("PG_SSL_MODE")

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

	return db
}
