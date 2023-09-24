package page_handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HomePageHandler struct{}

type CounterData struct {
	Value int
	Color string
}

var CounterColors = [10]string{"#fff", "#800", "#f00", "#080", "#0f0", "#008", "#00f", "#ff0", "#0ff", "#f0f"}

// Implements PageRouteRegistrar interface
func (h *HomePageHandler) Handler(engine *gin.Engine) {
	engine.GET("/", h.get)
	engine.POST("/component/home/counter", h.updateCounter)
}

func (h *HomePageHandler) get(c *gin.Context) {
	// Set additional page data
	data := gin.H{
		// Initial counter values
		"Counter": CounterData{
			Value: 0,
			Color: "#fff",
		},
	}

	// Render page
	c.HTML(http.StatusOK, "home_page.gohtml", data)
}

func (h *HomePageHandler) updateCounter(c *gin.Context) {
	currentCountStr := c.PostForm("currentCount")
	currentCount, err := strconv.Atoi(currentCountStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid currentCount value")
		return
	}

	// Update counter value
	currentCount++

	// Prepare updated counter data
	data := map[string]interface{}{
		// Initial counter values
		"Counter": CounterData{
			Value: currentCount,
			Color: CounterColors[currentCount%len(CounterColors)],
		},
	}

	// Render page
	c.HTML(http.StatusOK, "counter.gohtml", data)
}
