package config

import (
	"github.com/NotCoffee418/GoHtmxPgsql-Boilerplate/handlers/page_handlers"
	"github.com/NotCoffee418/GoHtmxPgsql-Boilerplate/internal/common"
)

// HandlerRegistrar Register all routes here, described in handlers
var RouteHandlers = []common.HandlerRegistrar{
	&page_handlers.HomeHandler{},
}
