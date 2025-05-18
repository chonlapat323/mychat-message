package main

import (
	"log"
	"mychat-message/database"
	"mychat-message/handlers"
	"mychat-message/middleware"
	"mychat-message/utils"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// โหลดค่าจาก .env
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()

		if err != nil {
			log.Fatal("ไม่พบไฟล์ .env หรือโหลดไม่สำเร็จ")
		}
	}

	database.InitMongo()
	utils.InitRedis()
	// /messages?room_id=xxx
	http.Handle("/messages", middleware.CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			middleware.JWTAuthMiddleware(http.HandlerFunc(handlers.GetMessagesHandler)).ServeHTTP(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	log.Println("Message service running on :5002")
	log.Fatal(http.ListenAndServe(":5002", nil))
}
