package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/NotCoffee418/GoHtmxPgsql-Boilerplate/config"
	"github.com/NotCoffee418/GoHtmxPgsql-Boilerplate/internal/server"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()

	// Register all routes here, described in handlers
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
