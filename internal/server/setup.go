package server

import (
	"embed"
	"fmt"
	"github.com/NotCoffee418/GoWebsite-Boilerplate/internal/config"
	"github.com/NotCoffee418/GoWebsite-Boilerplate/internal/handlers"
	"github.com/NotCoffee418/GoWebsite-Boilerplate/internal/types"
	"github.com/gin-contrib/gzip"
	"github.com/jmoiron/sqlx"
	"io/fs"
	"net/http"

	log "github.com/sirupsen/logrus"
	"os/exec"

	"github.com/gin-gonic/gin"
)

func SetupServer(engine *gin.Engine, db *sqlx.DB, templateFS embed.FS, staticFs embed.FS) {
	// Set up templates at templating.go
	initializeTemplates(engine, templateFS)

	// Background run postcss compile if enabled
	if config.DoMinifyCss {
		go runPostCSS(
			"./assets/dynamic/css/global.css",
			"./assets/dynamic/css/global.min.css")
	}

	// Set up /dynamic/ route
	engine.Static("/dynamic", "./assets/dynamic")

	// Set up /assets/static at / route
	registerStaticFiles(engine, staticFs)

	// Register middleware
	engine.Use(internalServerErrorHandlingMiddleware())
	engine.Use(gzip.Gzip(config.GzipCompressionLevel))

	// Prepare init context passed to all handlers.
	// Should include any global dependencies.
	handlerInitCtx := &types.HandlerInitContext{
		DB: db,
	}

	// Register all routes here
	for _, handler := range handlers.RouteHandlers {
		handler.Initialize(handlerInitCtx)
		handler.Handler(engine)
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

func registerStaticFiles(e *gin.Engine, staticFs embed.FS) {
	subFS, err := fs.Sub(staticFs, "assets/static")
	if err != nil {
		fmt.Println("Failed to create sub filesystem: ", err)
		return
	}

	walkDir := func(fsys fs.FS, fn func(path string, d fs.DirEntry, err error) error) error {
		return fs.WalkDir(fsys, ".", fn)
	}
	err = walkDir(subFS, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println("Error walking dir:", err)
			return err
		}
		if d.IsDir() {
			return nil
		}

		urlPath := "/" + filePath
		e.GET(urlPath, func(c *gin.Context) {
			c.FileFromFS(urlPath, http.FS(subFS))
		})

		return nil
	})
	if err != nil {
		log.Fatal("Error registering static files: ", err)
	}
}
