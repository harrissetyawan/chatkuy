package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/harrissetyawan/chatkuy/handlers"
)

func main() {
	viewsEngine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: viewsEngine,
	})

	app.Static("/static", "./static")
	appHandler := handlers.NewAppHandler()

	server := NewWebSocket()

	app.Get("/", appHandler.HandleGetIndex)
	app.Get("/ws", websocket.New(func(ctx *websocket.Conn) {
		server.HandleWebSocket(ctx)
	}))

	go server.HandleMessages()
	app.Listen(":3000")
}
