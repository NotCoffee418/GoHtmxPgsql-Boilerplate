package server

import (
	"github.com/jmoiron/sqlx"

	log "github.com/sirupsen/logrus"
	"os/exec"

	"github.com/NotCoffee418/GoWebsite-Boilerplate/config"
	"github.com/gin-gonic/gin"
)

func SetupServer(engine *gin.Engine, db *sqlx.DB) {
	// Set up templates at templating.go
	initializeTemplates(engine)

	// Background run postcss compile if enabled
	if config.DoMinifyCss {
		go runPostCSS("./static/css/global.css", "./static/css/global.min.css")
	}

	// Set up static file serving
	engine.Static("/static", "./static")

	// Register middleware
	engine.Use(internalServerErrorHandlingMiddleware())

	// Register all routes here
	for _, handler := range config.RouteHandlers {
		handler.Handler(engine, db)
	}

}

func runPostCSS(inputFile string, outputFile string) {
	log.Println("Running PostCSS...")
	cmd := exec.Command("npx", "postcss", inputFile, "-o", outputFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal("Error running postcss: ", string(output))
	}
	log.Println("PostCSS Ready:", string(output))
}
