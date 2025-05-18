package controller

import (
	"mychat-message/models"
	"mychat-message/service"

	"github.com/gofiber/fiber/v2"
)

func CreateMessage(c *fiber.Ctx) error {
	var msg models.Message

	// ✅ Parse body
	if err := c.BodyParser(&msg); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// ✅ Validate input
	if msg.RoomID == "" || msg.SenderID == "" || msg.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing fields"})
	}

	// ✅ Call service to save message
	if err := service.CreateMessage(&msg); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save message"})
	}

	return c.Status(fiber.StatusCreated).JSON(msg)
}

func GetMessages(c *fiber.Ctx) error {
	roomID := c.Params("roomId")
	if roomID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "roomId is required"})
	}

	messages, err := service.GetMessagesByRoom(roomID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get messages"})
	}

	return c.JSON(messages)
}
