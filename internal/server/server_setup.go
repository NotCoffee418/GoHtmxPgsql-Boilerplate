package server

import (
	"html/template"
	"log"
	"path/filepath"

	"github.com/NotCoffee418/GoHtmxPgsql-Boilerplate/config"
	"github.com/NotCoffee418/GoHtmxPgsql-Boilerplate/internal/utils"
	"github.com/gin-gonic/gin"
)

func SetupServer(engine *gin.Engine) {
	// Set up templates at common.Tmpl
	initializeTemplates(engine)

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
		func(path string) bool { return filepath.Ext(path) == ".html" })
	if err != nil {
		log.Fatal("Error listing templates: ", err)
	}
	tmpl, err := template.ParseFiles(all_templates...)
	if err != nil {
		log.Fatal("Error parsing templates: ", err)
	}
	engine.SetHTMLTemplate(tmpl)
}
