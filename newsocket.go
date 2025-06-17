package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

type WebSocketServer struct {
	id        string
	clients   map[*websocket.Conn]bool
	broadcast chan *Message
}

func NewWebSocket() *WebSocketServer {
	server := &WebSocketServer{
		id:        uuid.New().String(),
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan *Message),
	}
	log.Println("WebSocket ID:", server.clients)
	return server
}

func (s *WebSocketServer) HandleConnections(ctx *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(ctx) {
		return ctx.Next()
	}

	return fiber.ErrUpgradeRequired
}

func (s *WebSocketServer) HandleWebSocket(ctx *websocket.Conn) error {
	s.clients[ctx] = true
	log.Println("client connected:", ctx.RemoteAddr().String())
	defer func() {
		delete(s.clients, ctx)
		ctx.Close()
	}()

	for {
		_, msg, err := ctx.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
		}

		var message Message
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Fatalf("Error unmarshalling message: %v", err)
		}

		log.Printf("Message received: %+v", message)
		s.broadcast <- &message
	}

}

func (s *WebSocketServer) HandleMessages() {
	for {
		msg := <-s.broadcast

		for client := range s.clients {
			err := client.WriteMessage(websocket.TextMessage, getMessageTemplate(msg))
			if err != nil {
				log.Printf("Error sending message to client %s: %v", client.RemoteAddr().String(), err)
			}
		}

		log.Printf("Total active clients: %d", len(s.clients))
	}
}

func getMessageTemplate(msg *Message) []byte {
	tmpl, err := template.ParseFiles("views/messages.html")
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	var renderedMessage bytes.Buffer
	err = tmpl.Execute(&renderedMessage, msg)

	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}

	return renderedMessage.Bytes()
}
