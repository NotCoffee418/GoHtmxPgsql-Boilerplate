package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type HomeHandler struct{}

// Implements RouteRegistrar interface
func (h *HomeHandler) RegisterRoutes(router *mux.Router) {
	// GET Handler
	router.HandleFunc("/", h.Get).
		Methods("GET")

	// POST Handler
	router.HandleFunc("/", h.Post).
		Methods("POST")
}

func (h *HomeHandler) Get(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This Is Home")
}

func (h *HomeHandler) Post(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "POST on home")
}
