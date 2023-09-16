package handlers

import "net/http"

type HomeHandler struct{}

func (h *HomeHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", h.Handle)
}

func (h *HomeHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// Your code here
}
