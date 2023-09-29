package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/NotCoffee418/GoWebsite-Boilerplate/config"
	"github.com/NotCoffee418/GoWebsite-Boilerplate/db"
	"github.com/NotCoffee418/GoWebsite-Boilerplate/internal/page"
	"github.com/NotCoffee418/GoWebsite-Boilerplate/internal/server"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func main() {
	// Check if no arguments are provided
	if len(os.Args) == 1 {
		startServer()
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
			db.MigrateUp()
		case "down":
			db.MigrateDown()
		default:
			fmt.Println("Invalid migration command. Use 'migrate up' or 'migrate down'. For help, use the 'help' command.")
		}
	default:
		fmt.Println("Unknown command. For help, use the 'help' command.")
	}
}

func startServer() {
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

	// Background init database connection
	var conn *sqlx.DB
	go func() {
		conn = server.GetDBConn()
	}()
	defer func(conn *sqlx.DB) {
		_ = conn.Close()
	}(conn)

	// Set default page title when missing
	page.DefaultPageTitle = config.WebsiteTitle

	// Register all routes here, described in handlers
	engine := gin.Default()
	server.SetupServer(engine)

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
