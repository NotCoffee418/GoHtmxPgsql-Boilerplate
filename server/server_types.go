package server

import "net/http"

type RouteRegistrar interface {
	RegisterRoutes(mux *http.ServeMux)
}
