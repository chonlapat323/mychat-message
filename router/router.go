package router

import (
	"mychat-message/controller"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/messages", controller.CreateMessage)
	api.Get("/messages/:roomId", controller.GetMessages)
}
