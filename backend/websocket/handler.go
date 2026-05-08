package websocket

import (
	"log"
	"net/http"
	"whatsapp-clone/models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // In production, validate origin properly
	},
}

func HandleWebSocket(c *gin.Context) {
	userID := c.GetUint("user_id")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	client := &Client{
		UserID: userID,
		conn:   conn,
		send:   make(chan models.Message, 256),
	}

	hub.register <- client

	go client.writePump()
	go client.readPump()
}
