package homepage

import (
	"github.com/NotCoffee418/GoWebsite-Boilerplate/internal/types"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler Implements types.HandlerRegistrar interface
type Handler struct{}

type CounterData struct {
	Value int
	Color string
}

var (
	CounterColors = [10]string{"#fff", "#800", "#f00", "#080", "#0f0", "#008", "#00f", "#ff0", "#0ff", "#f0f"}
)

// Initialize is called before the handler is registered
func (h *Handler) Initialize(_ *types.HandlerInitContext) {
	// Nothing to initialize
}

// Handler Implements RouteRegistrar interface
func (h *Handler) Handler(engine *gin.Engine) {
	engine.GET("/", h.get)
	engine.POST("/component/home/counter", h.updateCounter)
}

func (h *Handler) get(c *gin.Context) {
	// Set SEO meta data
	meta := &types.MetaData{
		Title:       "Demo Home Page",
		Description: "This is a demo home page showing off the boilerplate.",
	}

	// Initial counter values
	data := &map[string]interface{}{
		"Counter": CounterData{
			Value: 0,
			Color: "#fff",
		},
	}

	structuredData := types.StructurePageData(&data, meta)

	// Render page
	c.HTML(http.StatusOK, "home_page.html", structuredData)
}

func (h *Handler) updateCounter(c *gin.Context) {
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

	structuredData := types.StructurePageData(&data, nil)

	// Render page
	c.HTML(http.StatusOK, "counter.html", structuredData)
}
