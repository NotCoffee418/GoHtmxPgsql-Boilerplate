package main

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"time"

	"github.com/NotCoffee418/GoWebsite-Boilerplate/config"
	"github.com/NotCoffee418/GoWebsite-Boilerplate/internal/page"
	"github.com/NotCoffee418/GoWebsite-Boilerplate/internal/server"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

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
			<-server.MigrateUpCh(db)
		case "down":
			<-server.MigrateDownCh(db)
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
	liveState := <-server.GetLiveMigrationInfoCh(db)
	if liveState.InstalledVersion < liveState.AvailableVersion {
		log.Warnf("Database migration required. Currently at %d, available version is %d.",
			liveState.InstalledVersion, liveState.AvailableVersion)
	}

	// Set default page title when missing
	page.DefaultPageTitle = config.WebsiteTitle

	// Register all routes here, described in handlers
	engine := gin.Default()
	server.SetupServer(engine, db)

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
