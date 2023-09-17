package internal

import (
	"github.com/NotCoffee418/GoHtmxPgsql-Boilerplate/config"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterRoutes(router *mux.Router) {
	// Redirect to HTTPS
	if config.HttpsRedirect {
		router.HandleFunc("/", redirectToHttps).Schemes("http")
	}

	// Register all routes here, described in handlers
	for _, registrar := range *config.RouteHandlers {
		registrar.RegisterRoutes(router)
	}
}

func redirectToHttps(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+r.Host+r.URL.Path, http.StatusMovedPermanently)
}
