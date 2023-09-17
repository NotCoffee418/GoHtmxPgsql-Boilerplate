package server_models

import "github.com/gorilla/mux"

type RouteRegistrar interface {
	RegisterRoutes(router *mux.Router)
}
