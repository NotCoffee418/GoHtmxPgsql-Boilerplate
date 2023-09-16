package main

import (
	"github.com/NotCoffee418/GoHtmxPgsql-Boilerplate/handlers"
	"github.com/NotCoffee418/GoHtmxPgsql-Boilerplate/server"
)

const WebsiteTitle = ""

// Register all routes here, described in handlers
var route_registrars = &[]server.RouteRegistrar{
	&handlers.HomeHandler{},
}
