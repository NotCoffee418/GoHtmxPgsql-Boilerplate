package homepage

import (
	"github.com/NotCoffee418/GoWebsite-Boilerplate/internal/types"
	"github.com/NotCoffee418/GoWebsite-Boilerplate/pkg/guestbook"
	"github.com/NotCoffee418/websocketmanager"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
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
var wsManager = websocketmanager.NewDefaultManager()
var db *sqlx.DB
var gbRepo = guestbook.Repository{}

// Handler Implements PageRouteRegistrar interface
func (h *HomePageHandler) Handler(engine *gin.Engine, _db *sqlx.DB) {
	db = _db
	engine.GET("/", h.get)
	engine.POST("/component/home/counter", h.updateCounter)
	engine.GET("/home/guestbook/ws", h.guestbookWS)
	engine.POST("/home/guestbook/submit", h.guestbookPost)
	engine.GET("/home/guestbook/recent", h.getRecentGuestbookEntries)
}

func (h *HomePageHandler) get(c *gin.Context) {
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

	structuredData := types.StructurePageData(&data, nil)

	// Render page
	c.HTML(http.StatusOK, "counter.html", structuredData)
}

func (h *HomePageHandler) guestbookWS(context *gin.Context) {
	wsManager.UpgradeClientCh(context.Writer, context.Request)
}

func (h *HomePageHandler) guestbookPost(context *gin.Context) {
	// Get form data
	name := context.PostForm("name")
	message := context.PostForm("message")

	// Validation logic
	if len(name) == 0 {
		context.String(http.StatusBadRequest, "Name cannot be empty")
		return
	}
	if len(name) > 255 {
		context.String(http.StatusBadRequest, "Name cannot be longer than 50 characters")
		return
	}
	if len(message) == 0 {
		context.String(http.StatusBadRequest, "Message cannot be empty")
		return
	}

	// Insert into database
	err := gbRepo.Insert(db, name, message)
	if err != nil {
		context.String(http.StatusInternalServerError, "Error inserting into database")
		return
	}

	// Broadcast update to all clients
	h.triggerGuestbookUpdate()
}

func (h *HomePageHandler) triggerGuestbookUpdate() {
	wsManager.BroadcastMessage(websocket.TextMessage, []byte("New guestbook entry"))
}

func (h *HomePageHandler) getRecentGuestbookEntries(c *gin.Context) {
	//c.
}
