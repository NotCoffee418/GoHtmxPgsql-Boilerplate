package handlers

import (
	"github.com/NotCoffee418/GoWebsite-Boilerplate/internal/types"
	"github.com/NotCoffee418/GoWebsite-Boilerplate/pkg/guestbook"
	"github.com/NotCoffee418/GoWebsite-Boilerplate/pkg/homepage"
)

// RouteHandlers Register all routes here, described in handlers
var RouteHandlers = []types.HandlerRegistrar{
	&homepage.HomeApiHandler{},
	&homepage.HomePageHandler{},
	&guestbook.GuestbookHandler{},
}
