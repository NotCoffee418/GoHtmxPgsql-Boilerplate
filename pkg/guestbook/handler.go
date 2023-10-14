package guestbook

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/NotCoffee418/GoWebsite-Boilerplate/internal/config"
	"github.com/NotCoffee418/GoWebsite-Boilerplate/internal/types"
	"github.com/NotCoffee418/GoWebsite-Boilerplate/internal/utils"
	"github.com/NotCoffee418/GoWebsite-Boilerplate/internal/weberrors"
	"github.com/NotCoffee418/GoWebsite-Boilerplate/pkg/form"
	"github.com/NotCoffee418/websocketmanager"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
)

var (
	wsManager *websocketmanager.Manager
	db        *sqlx.DB
	gbRepo    = Repository{}
)

// Handler Implements types.HandlerRegistrar interface
type GuestbookHandler struct{}

// Initialize is called before the handler is registered
func (h *GuestbookHandler) Initialize(initContext *types.HandlerInitContext) {
	db = initContext.DB

	// Instance for the websocket manager is tied to the component
	wsManager = websocketmanager.NewBuilder().
		WithReadBufferSize(config.WsReadBufferSize).
		WithWriteBufferSize(config.WsWriteBufferSize).
		Build()
}

// Handler Handler Implements RouteRegistrar interface
func (h *GuestbookHandler) Handler(engine *gin.Engine) {
	engine.GET("/home/guestbook/ws", h.guestbookWsInvoker)
	engine.GET("/home/guestbook/recent", h.getRecentGuestbookEntries)
}

func (h *GuestbookHandler) guestbookWsInvoker(context *gin.Context) {
	client, err := wsManager.UpgradeClient(context.Writer, context.Request)
	if err != nil {
		log.Errorf("Failed to upgrade client: %v", err)
		return
	}
	wsManager.RegisterClientObserver(*client.ConnId, h.guestbookWsHandler)
}

func (h *GuestbookHandler) guestbookWsHandler(wsClient *websocketmanager.Client, messageType int, wsData []byte) {
	// Get submission data
	success := true
	responseBuilder := form.NewFormResponseBuilder("guestbook-form-validation").
		SetContainerClasses("mb-4")

	// Parse incoming message
	wsMsg, err := utils.DeserializeHtmxWebsocketMessage(string(wsData))
	if err != nil {
		// This message is not for this function
		return
	}
	name, name_ok := wsMsg.Data["name"]
	message, message_ok := wsMsg.Data["message"]
	if !name_ok || !message_ok {
		success = false
		responseBuilder.AddMessage(false, "400: Bad Request")
	}

	// Validation logic
	if success {
		if len(name) == 0 {
			responseBuilder.AddMessage(false, "Name cannot be empty")
			success = false
		}
		if len(name) > 255 {
			responseBuilder.AddMessage(false, "Name cannot be longer than 50 characters")
			success = false
		}
		if len(message) < 4 {
			responseBuilder.AddMessage(false, "Message is too short")
			success = false
		}
	}

	// Insert into database
	if success {
		err := gbRepo.Insert(db, name, message)
		if err != nil {
			responseBuilder.AddMessage(false, "Internal server error")
			success = false
		}
	}

	// Valdiation message
	if success {
		// All is well, add success message
		responseBuilder.AddMessage(true, "Message Posted!")

		// Broadcast update to all clients
		h.triggerGuestbookUpdate()
	}

	// Prepare ws response
	validationHtml, err := responseBuilder.BuildHtmlString()
	if err != nil {
		log.Errorf("Failed to get validation html: %v", err)
	}

	// Send it
	wsManager.SendMessageToClient(*wsClient.ConnId, websocket.TextMessage, []byte(validationHtml))

}

func (h *GuestbookHandler) triggerGuestbookUpdate() {
	recentPosts, err := gbRepo.GetRecent(db, 10)
	if err != nil {
		return
	}
	data := map[string]interface{}{
		"Posts": NewMessageDisplaySlice(recentPosts),
	}
	html, err := utils.RenderedTemplateString("guestbook_posts.html", data)
	html = utils.OobSwapWrap("guestbook-messages", html)
	if err != nil {
		log.Errorf("Failed to render guestbook posts: %v", err)
		return
	}
	wsManager.BroadcastMessage(websocket.TextMessage, []byte(html))
}

func (h *GuestbookHandler) getRecentGuestbookEntries(c *gin.Context) {
	recentPosts, err := gbRepo.GetRecent(db, 10)
	if err != nil {
		weberrors.InternalServerErrorResponse(c, err)
		return
	}

	data := map[string]interface{}{
		"Posts": NewMessageDisplaySlice(recentPosts),
	}

	c.HTML(http.StatusOK, "guestbook_posts.html", &data)
}
