package homepage

import (
	"net/http"
	"strconv"
	"time"

	"github.com/NotCoffee418/GoWebsite-Boilerplate/internal/types"

	"github.com/gin-gonic/gin"
)

// Handler Implements types.HandlerRegistrar interface
type HomePageHandler struct{}

type CounterData struct {
	Value int
	Color string
}

var (
	CounterColors = [10]string{"#fff", "#800", "#f00", "#080", "#0f0", "#008", "#00f", "#ff0", "#0ff", "#f0f"}
)

// Initialize is called before the handler is registered
func (h *HomePageHandler) Initialize(_ *types.HandlerInitContext) {
	// Nothing to initialize
}

// Handler Implements RouteRegistrar interface
func (h *HomePageHandler) Handler(engine *gin.Engine) {
	engine.GET("/", h.get)
	engine.POST("/component/home/counter", h.updateCounter)
}

func (h *HomePageHandler) get(c *gin.Context) {
	// Set SEO meta data
	meta := &types.MetaData{
		Title:       "Demo Home Page",
		Description: "This is a demo home page showing off the boilerplate.",
	}

	// Initial counter values
	data := map[string]interface{}{
		"Counter": CounterData{
			Value: 0,
			Color: "#fff",
		},
	}

	structuredData := types.NewStructurePageData(data, meta)

	// Render page
	c.HTML(http.StatusOK, "home_page.html", structuredData)
}

func (h *HomePageHandler) updateCounter(c *gin.Context) {
	time.Sleep(1 * time.Second) // Artificial delay to demo htmx loading indicator
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

	structuredData := types.NewStructurePageData(data, nil)

	// Render page
	c.HTML(http.StatusOK, "counter.html", structuredData)
}
