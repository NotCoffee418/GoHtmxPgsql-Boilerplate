package guestbook

import (
	"github.com/NotCoffee418/GoWebsite-Boilerplate/internal/config"
	"github.com/NotCoffee418/GoWebsite-Boilerplate/internal/types"
	"github.com/NotCoffee418/GoWebsite-Boilerplate/internal/weberrors"
	"github.com/NotCoffee418/websocketmanager"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"net/http"
)

var (
	wsManager *websocketmanager.Manager
	db        *sqlx.DB
	gbRepo    = Repository{}
)

// Handler Implements types.HandlerRegistrar interface
type Handler struct{}

// Initialize is called before the handler is registered
func (h *Handler) Initialize(initContext *types.HandlerInitContext) {
	db = initContext.DB

	// Instance for the websocket manager is tied to the component
	wsManager = websocketmanager.NewBuilder().
		WithReadBufferSize(config.WsReadBufferSize).
		WithWriteBufferSize(config.WsWriteBufferSize).
		Build()
}

// GuestbookHandler Handler Implements RouteRegistrar interface
func (h *Handler) GuestbookHandler(engine *gin.Engine) {
	engine.GET("/home/guestbook/ws", h.guestbookWS)
	engine.POST("/home/guestbook/submit", h.guestbookPost)
	engine.GET("/home/guestbook/recent", h.getRecentGuestbookEntries)
}

func (h *Handler) guestbookWS(context *gin.Context) {
	<-wsManager.UpgradeClientCh(context.Writer, context.Request)
}

func (h *Handler) guestbookPost(context *gin.Context) {
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

func (h *Handler) triggerGuestbookUpdate() {
	wsManager.BroadcastMessage(websocket.TextMessage, []byte("New guestbook entry"))
}

func (h *Handler) getRecentGuestbookEntries(c *gin.Context) {
	recent, err := gbRepo.GetRecent(db, 10)
	if err != nil {
		weberrors.InternalServerErrorResponse(c, err)
		return
	}

}
