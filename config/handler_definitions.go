package config

import (
	"github.com/NotCoffee418/GoHtmxPgsql-Boilerplate/handlers"
	"github.com/NotCoffee418/GoHtmxPgsql-Boilerplate/internal/server_models"
)

// RouteHandlers Register all routes here, described in handlers
var RouteHandlers = &[]server_models.RouteRegistrar{
	&handlers.HomeHandler{},
}
