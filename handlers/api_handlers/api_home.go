package api_handlers

import (
	"github.com/jmoiron/sqlx"
	"net/http"
	"time"

	"github.com/NotCoffee418/GoWebsite-Boilerplate/internal/common"
	"github.com/gin-gonic/gin"
)

type HomeApiHandler struct{}

type HomePageData struct {
	Time string `json:"time"`
}

// Handler Implements PageRouteRegistrar interface
func (h *HomeApiHandler) Handler(engine *gin.Engine, _ *sqlx.DB) {
	engine.GET("/api/home/get-server-time", h.get)
}

func (h *HomeApiHandler) get(c *gin.Context) {
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	resp := common.ApiResponseFactory.Ok(
		&HomePageData{Time: timeStr})

	// Render page
	c.JSON(http.StatusOK, resp)
}
