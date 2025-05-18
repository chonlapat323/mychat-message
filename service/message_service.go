package service

import (
	"mychat-message/models"
	"mychat-message/repository"
)

func CreateMessage(message *models.Message) error {
	return repository.CreateMessage(message)
}

func GetMessagesByRoom(roomID string) ([]models.Message, error) {
	return repository.GetMessagesByRoom(roomID)
}
