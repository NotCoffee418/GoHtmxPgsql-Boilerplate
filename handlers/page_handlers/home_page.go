package page_handlers

import (
	"github.com/NotCoffee418/GoHtmxPgsql-Boilerplate/internal/common"
	"github.com/gin-gonic/gin"
)

type HomeHandler struct{}

// Implements PageRouteRegistrar interface
func (h *HomeHandler) Handler(engine *gin.Engine) {
	engine.GET("/", h.get)
}

func (h *HomeHandler) get(c *gin.Context) {
	if err := common.Tmpl.ExecuteTemplate(c.Writer, "home.html", nil); err != nil {
		c.Error(err).SetType(gin.ErrorTypePrivate) // This sets the error type. You can handle it in a centralized middleware.
	}
}
