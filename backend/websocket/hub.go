package websocket

import (
	"whatsapp-clone/models"
)

type Hub struct {
	clients    map[uint]*Client // userID -> Client
	broadcast  chan models.Message
	register   chan *Client
	unregister chan *Client
}

var hub = &Hub{
	clients:    make(map[uint]*Client),
	broadcast:  make(chan models.Message),
	register:   make(chan *Client),
	unregister: make(chan *Client),
}

func GetHub() *Hub {
	return hub
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.UserID] = client

		case client := <-h.unregister:
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.send)
			}

		case message := <-h.broadcast:
			// Send to all clients in the conversation
			// Note: In a production app, you'd want to fetch conversation members
			// and only send to them. For simplicity, we're broadcasting to all.
			for _, client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client.UserID)
				}
			}
		}
	}
}

func BroadcastMessage(message models.Message) {
	hub.broadcast <- message
}

func init() {
	go hub.Run()
}
