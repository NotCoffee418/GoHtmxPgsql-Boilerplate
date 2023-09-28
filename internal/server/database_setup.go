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
	PG_HOST := common.GetEnv("PG_HOST")
	PG_PORT := common.GetEnv("PG_PORT")
	PG_USER := common.GetEnv("PG_USER")
	PG_PASS := common.GetEnv("PG_PASS")
	PG_DATABASE := common.GetEnv("PG_DATABASE")
	PG_SSL_MODE := common.GetEnv("PG_SSL_MODE")

	// Validate SSL mode
	validSslModes := []string{
		"disable", "allow", "prefer", "require", "verify-ca", "verify-full"}
	if !utils.SliceContains(validSslModes, PG_SSL_MODE) {
		log.Fatalf("Invalid PG_SSL_MODE: %s", PG_SSL_MODE)
	}

	// Set up connection
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		PG_HOST, PG_PORT, PG_USER, PG_PASS, PG_DATABASE, PG_SSL_MODE)
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
