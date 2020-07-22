package main

import (
	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"
	"github.com/gofiber/helmet"
	"github.com/gofiber/logger"
	"github.com/gofiber/websocket"
	"log"
)

func main() {
	app := fiber.New()
	app.Use(helmet.New())
	app.Use(logger.New())
	app.Use(cors.New())

	// create hub for coordination between clients
	hub := NewHub()

	// websocket connection handler
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		// create new client with current connection
		NewClient(hub, c)
	}))

	// just listen
	err := app.Listen(4500)
	if err != nil {
		log.Fatalln(err)
	}
}
