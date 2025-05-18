package handlers

import (
	"context"
	"encoding/json"
	"log"
	"mychat-message/database"
	"mychat-message/models"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("🟡 [Handler] GetMessagesHandler called")
	roomIDStr := r.URL.Query().Get("room_id")
	log.Println("🔍 GET /messages?room_id=", roomIDStr)
	if roomIDStr == "" {
		http.Error(w, "Missing room_id", http.StatusBadRequest)
		return
	}

	roomID, err := primitive.ObjectIDFromHex(roomIDStr)

	log.Println("➡️ Querying Mongo with ObjectID:", roomID)
	if err != nil {
		http.Error(w, "Invalid room_id", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := database.MessageCollection.Find(ctx, bson.M{"room_id": roomID})
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var messages []models.Message
	if err = cursor.All(ctx, &messages); err != nil {
		http.Error(w, "Failed to decode messages", http.StatusInternalServerError)
		return
	}

	writeJSON(w, messages)
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
