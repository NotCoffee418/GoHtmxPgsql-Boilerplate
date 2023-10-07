package main

import (
	"embed"
	"errors"
	"fmt"
	"github.com/NotCoffee418/GoWebsite-Boilerplate/internal/config"
	"github.com/NotCoffee418/GoWebsite-Boilerplate/internal/types"
	"github.com/NotCoffee418/GoWebsite-Boilerplate/internal/utils"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"time"

	"github.com/NotCoffee418/GoWebsite-Boilerplate/internal/server"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

//go:embed all:migrations
var migrationFS embed.FS

//go:embed all:templates
var templatesFS embed.FS

//go:embed all:assets/static
var staticFs embed.FS

func main() {
	// Enable optional .env file
	_, err := os.Stat(".env")
	if errors.Is(err, os.ErrNotExist) {
		log.Println("No .env file found, skipping...")
	} else {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file:", err)
		} else {
			log.Println("Loaded .env file")
		}
	}

	// init database connection
	log.Println("Connecting to database...")
	db := server.GetDB()
	defer func(conn *sqlx.DB) {
		_ = conn.Close()
	}(db)

	// Check if no arguments are provided
	if len(os.Args) == 1 {
		startServer(db)
		return
	}

	// Parse command arguments
	switch os.Args[1] {
	case "help":
		showHelp()
	case "migrate":
		if len(os.Args) < 3 {
			fmt.Println("Expected 'up' or 'down' argument after 'migrate'. For help, use the 'help' command.")
			return
		}
		switch os.Args[2] {
		case "up":
			<-server.MigrateUpCh(db, migrationFS)
		case "down":
			<-server.MigrateDownCh(db, migrationFS)
		default:
			fmt.Println("Invalid migration command. Use 'migrate up' or 'migrate down'. For help, use the 'help' command.")
		}
	default:
		fmt.Println("Unknown command. For help, use the 'help' command.")
	}
	log.Println("Clean exit.")
}

func startServer(db *sqlx.DB) {
	// Database migration check
	log.Println("Checking database migration status...")
	liveState := <-server.GetLiveMigrationInfoCh(db, migrationFS)
	if liveState.InstalledVersion < liveState.AvailableVersion {
		if utils.GetEnv("AUTO_APPLY_MIGRATIONS") == "true" {
			log.Println("Migrating database...")
			<-server.MigrateUpCh(db, migrationFS)
			log.Println("Database migration complete.")
		} else {
			log.Warn("Database migration required. Set MIGRATE_ON_START to true to automatically migrate.")
		}
	}

	// Set default page title when missing
	types.DefaultPageTitle = config.WebsiteTitle

	// Set up routes, middleware, etc.
	engine := gin.Default()
	server.SetupServer(engine, db, templatesFS, staticFs)

	// Start server
	svr := &http.Server{
		Handler:      engine,
		Addr:         fmt.Sprintf("127.0.0.1:%d", config.ListenPort),
		WriteTimeout: config.TimeoutSeconds * time.Second,
		ReadTimeout:  config.TimeoutSeconds * time.Second,
	}
	log.Print("Server started on port ", config.ListenPort)
	log.Fatal(svr.ListenAndServe())
}

func showHelp() {
	fmt.Println(`Available commands:
	help           - Show this help message.
	migrate up     - Apply migrations.
	migrate down   - Rollback migrations.`)
}
