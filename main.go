package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/NotCoffee418/GoHtmxPgsql-Boilerplate/config"
	"github.com/NotCoffee418/GoHtmxPgsql-Boilerplate/internal/server"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
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

	// Background init datbase connection
	var conn *sqlx.DB
	go func() {
		conn = server.GetConn()
	}()

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

	// Close DB connection
	if conn != nil {
		log.Println("Closing DB connection...")
		defer conn.Close()
	}
}
