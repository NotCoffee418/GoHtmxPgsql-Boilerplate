package api_handlers

import (
	"github.com/NotCoffee418/GoWebsite-Boilerplate/internal/common"
)

// RouteHandlers Register all routes here, described in handlers
var RouteHandlers = []common.HandlerRegistrar{
	&HomeApiHandler{},
}
