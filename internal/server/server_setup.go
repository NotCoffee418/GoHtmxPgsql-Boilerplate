package server

import (
	"html/template"
	"log"
	"os/exec"
	"path/filepath"

	"github.com/NotCoffee418/GoHtmxPgsql-Boilerplate/config"
	"github.com/NotCoffee418/GoHtmxPgsql-Boilerplate/internal/utils"
	"github.com/gin-gonic/gin"
)

func SetupServer(engine *gin.Engine) {
	// Set up templates at common.Tmpl
	initializeTemplates(engine)

	// Background run postcss if enabled
	if config.DoMinifyCss {
		go runPostCSS()
	}

	// Set up static file serving
	engine.Static("/static", "./static")

	// Register middleware
	engine.Use(internalServerErrorHandlingMiddleware())

	// Register all routes here
	for _, handler := range config.RouteHandlers {
		handler.Handler(engine)
	}

}

func initializeTemplates(engine *gin.Engine) {
	// Load HTML templates
	all_templates, err := utils.GetRecursiveFiles(
		"./templates",
		func(path string) bool { return filepath.Ext(path) == ".gohtml" })
	if err != nil {
		log.Fatal("Error listing templates: ", err)
	}
	tmpl, err := template.ParseFiles(all_templates...)
	if err != nil {
		log.Fatal("Error parsing templates: ", err)
	}
	engine.SetHTMLTemplate(tmpl)
}

func runPostCSS() {
	cmd := exec.Command(
		"npx", "postcss", "./static/css/global.css",
		"-o", "./static/css/global.min.css")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal("Error running postcss: ", string(output))
	}
	log.Println("PostCSS Ready:", string(output))
}
