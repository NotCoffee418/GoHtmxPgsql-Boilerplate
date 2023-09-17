package main

import (
	"fmt"
	"github.com/NotCoffee418/GoHtmxPgsql-Boilerplate/config"
	"github.com/NotCoffee418/GoHtmxPgsql-Boilerplate/internal"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func main() {
	router := mux.NewRouter()

	// Register all routes here, described in handlers
	internal.RegisterRoutes(router)

	// Serve static files
	router.PathPrefix("/").
		Handler(http.StripPrefix("/", http.FileServer(http.Dir("static/"))))

	// Start server
	svr := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("127.0.0.1:%d", config.ListenPort),
		WriteTimeout: config.TimeoutSeconds * time.Second,
		ReadTimeout:  config.TimeoutSeconds * time.Second,
	}
	log.Print("Server started on port ", config.ListenPort)
	log.Fatal(svr.ListenAndServe())
}
