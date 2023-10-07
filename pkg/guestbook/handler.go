package guestbook

import (
	"github.com/NotCoffee418/GoWebsite-Boilerplate/internal/config"
	"github.com/NotCoffee418/GoWebsite-Boilerplate/internal/types"
	"github.com/NotCoffee418/GoWebsite-Boilerplate/internal/weberrors"
	"github.com/NotCoffee418/GoWebsite-Boilerplate/pkg/validation"
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
	engine.GET("/home/guestbook/ws", h.guestbookWsInvoker)
	engine.POST("/home/guestbook/submit", h.guestbookPost)
	engine.GET("/home/guestbook/recent", h.getRecentGuestbookEntries)
}

func (h *Handler) guestbookWsInvoker(context *gin.Context) {
	client := <-wsManager.UpgradeClientCh(context.Writer, context.Request)
	wsManager.RegisterClientObserver(*client.ConnId, guestbookWsHandler)
}

func guestbookWsHandler(wsClient *websocketmanager.Client, messageType int, message []byte) {
	wsManager.BroadcastMessage(websocket.TextMessage, []byte("New guestbook entry"))
}

func (h *Handler) guestbookPost(c *gin.Context) {
	// Get form data
	name := c.PostForm("name")
	message := c.PostForm("message")

	success := true
	var respMsgs []string

	// Validation logic
	if len(name) == 0 {
		respMsgs = append(respMsgs, "Name cannot be empty")
		success = false
	}
	if len(name) > 255 {
		respMsgs = append(respMsgs, "Name cannot be longer than 50 characters")
		success = false
	}
	if len(message) < 4 {
		respMsgs = append(respMsgs, "Message cannot be empty")
		success = false
	}

	// Insert into database
	if success {
		err := gbRepo.Insert(db, name, message)
		if err != nil {
			respMsgs = append(respMsgs, "Internal server error")
			success = false
		}
	}

	if success {
		// All is well, add success message
		respMsgs = append(respMsgs, "Successfully added post")

		// Broadcast update to all clients
		h.triggerGuestbookUpdate()
	}

	// Render form response
	var respData []validation.FormResponse
	for _, msg := range respMsgs {
		respData = append(respData, validation.FormResponse{
			IsSuccessful: success,
			Message:      msg,
		})
	}
	validation.RenderFormResponse(c, respData)
}

func (h *Handler) triggerGuestbookUpdate() {
	go wsManager.BroadcastMessage(websocket.TextMessage, []byte("New guestbook entry"))
}

func (h *Handler) getRecentGuestbookEntries(c *gin.Context) {
	recentPosts, err := gbRepo.GetRecent(db, 10)
	if err != nil {
		weberrors.InternalServerErrorResponse(c, err)
		return
	}

	data := &map[string]interface{}{
		"Posts": recentPosts,
	}

	c.HTML(http.StatusOK, "guestbook_posts.html", data)
}
