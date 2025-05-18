package repository

import (
	"context"
	"fmt"
	"time"

	"mychat-message/database"
	"mychat-message/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func getCollection() *mongo.Collection {
	return database.MessageCollection
}

func CreateMessage(message *models.Message) error {
	message.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	result, err := getCollection().InsertOne(context.TODO(), message)
	if err != nil {
		fmt.Println("❌ Insert error:", err)
		return err
	}
	fmt.Println("✅ Inserted ID:", result.InsertedID)
	return err
}

func GetMessagesByRoom(roomID string) ([]models.Message, error) {
	var messages []models.Message

	cursor, err := getCollection().Find(context.TODO(), bson.M{"room_id": roomID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var msg models.Message
		if err := cursor.Decode(&msg); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
