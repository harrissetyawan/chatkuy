package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/harrissetyawan/chatkuy/handlers"
)

func main() {
	viewsEngine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: viewsEngine,
	})

	app.Static("/static", "./static")

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello, World!")
	})
	app.Listen(":3000")
}
