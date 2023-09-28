package config

import (
	"github.com/NotCoffee418/GoWebsite-Boilerplate/handlers/api_handlers"
	"github.com/NotCoffee418/GoWebsite-Boilerplate/handlers/page_handlers"
	"github.com/NotCoffee418/GoWebsite-Boilerplate/internal/common"
)

// HandlerRegistrar Register all routes here, described in handlers
var RouteHandlers = []common.HandlerRegistrar{
	&page_handlers.HomePageHandler{},
	&api_handlers.HomeApiHandler{},
}
