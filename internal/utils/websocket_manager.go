package utils

import (
	"github.com/NotCoffee418/GoWebsite-Boilerplate/config"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
	"time"
)

type WebSocketManager struct {
	clients map[*websocket.Conn]uuid.UUID
	mutex   *sync.Mutex
}

func NewWebSocketManager() *WebSocketManager {
	wm := &WebSocketManager{
		clients: make(map[*websocket.Conn]uuid.UUID),
		mutex:   &sync.Mutex{},
	}
	go wm.cleanupClosedClients()
	return wm
}

// UpgradeClient upgrades a client connection to a websocket connection
// Returns a channel that will receive the UUID of the client
func (wm *WebSocketManager) UpgradeClient(c *gin.Context) chan uuid.UUID {
	upgradeChan := make(chan uuid.UUID, 1)
	go func() {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  config.WsReadBufferSize,
			WriteBufferSize: config.WsWriteBufferSize,
		}
		wsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error creating websocket connection")
			return
		}
		upgradeChan <- wm.Register(wsConn)
		close(upgradeChan)
	}()
	return upgradeChan
}

func (wm *WebSocketManager) Register(conn *websocket.Conn) uuid.UUID {
	userID := uuid.New()
	wm.mutex.Lock()
	wm.clients[conn] = userID
	wm.mutex.Unlock()
	return userID
}

func (wm *WebSocketManager) Unregister(conn *websocket.Conn) {
	wm.mutex.Lock()
	delete(wm.clients, conn)
	wm.mutex.Unlock()
}

func (wm *WebSocketManager) BroadcastMessage(message []byte) {
	wm.mutex.Lock()
	defer wm.mutex.Unlock()

	for client := range wm.clients {
		if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("Websocket error: %s", err)
			_ = client.Close()
			delete(wm.clients, client)
		}
	}
}

func (wm *WebSocketManager) SendMessageToUser(userID uuid.UUID, message []byte) {
	wm.mutex.Lock()
	defer wm.mutex.Unlock()

	for client, id := range wm.clients {
		if id == userID {
			if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("Websocket error: %s", err)
				_ = client.Close()
				delete(wm.clients, client)
			}
		}
	}
}

func (wm *WebSocketManager) cleanupClosedClients() {
	ticker := time.NewTicker(30 * time.Second) // Ping every 30 seconds

	for {
		select {
		case <-ticker.C:
			wm.pingClients()
		}
	}
}

func (wm *WebSocketManager) pingClients() {
	wm.mutex.Lock()
	defer wm.mutex.Unlock()

	for client := range wm.clients {
		if err := client.WriteMessage(websocket.PingMessage, nil); err != nil {
			log.Printf("Websocket error: %s", err)
			client.Close()
			delete(wm.clients, client)
		}
	}
}
