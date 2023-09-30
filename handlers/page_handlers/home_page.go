package page_handlers

import (
	"fmt"
	"github.com/NotCoffee418/GoWebsite-Boilerplate/internal/server"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"

	"github.com/NotCoffee418/GoWebsite-Boilerplate/internal/page"
	"github.com/gin-gonic/gin"
)

type HomePageHandler struct{}

type CounterData struct {
	Value int
	Color string
}

var CounterColors = [10]string{"#fff", "#800", "#f00", "#080", "#0f0", "#008", "#00f", "#ff0", "#0ff", "#f0f"}
var wsManager = server.NewWebSocketManager()

// Handler Implements PageRouteRegistrar interface
func (h *HomePageHandler) Handler(engine *gin.Engine, _ *sqlx.DB) {
	engine.GET("/", h.get)
	engine.POST("/component/home/counter", h.updateCounter)
	engine.GET("/home/guestbook/ws", h.guestbookWS)
}

func (h *HomePageHandler) get(c *gin.Context) {
	// Set SEO meta data
	meta := &page.MetaData{
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

	structuredData := page.StructurePageData(&data, meta)

	// Render page
	c.HTML(http.StatusOK, "home_page.html", structuredData)
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

	structuredData := page.StructurePageData(&data, nil)

	// Render page
	c.HTML(http.StatusOK, "counter.html", structuredData)
}

func (h *HomePageHandler) guestbookWS(c *gin.Context) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	wsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating websocket connection")
		return
	}
	if err != nil {
		fmt.Println("Error generating UUID:", err)
		return
	}
	wsManager.Register(wsConn)
}

func (h *HomePageHandler) triggerGuestbookUpdate() {

}
