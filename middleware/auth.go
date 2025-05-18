package middleware

import (
	"context"
	"log"
	"mychat-message/contextkey"
	"mychat-message/utils"
	"net/http"
)

type contextKey string

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("🔑 JWTAuthMiddleware called")

		defer func() {
			if rec := recover(); rec != nil {
				log.Println("🔥 Panic in JWT middleware:", rec)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		cookie, err := r.Cookie("token")
		if err != nil || cookie.Value == "" {
			log.Println("❌ Missing or invalid token:", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := cookie.Value
		log.Println("🧪 Token received:", tokenString[:10], "...")

		isBlacklisted, err := utils.IsTokenBlacklisted(tokenString)
		if err != nil {
			log.Println("❌ Redis error or client nil:", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		if isBlacklisted {
			log.Println("🚫 Token is blacklisted")
			http.Error(w, "Token revoked", http.StatusUnauthorized)
			return
		}

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			log.Println("❌ JWT validation failed:", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		log.Println("✅ JWT valid for user:", claims.UserID)
		ctx := context.WithValue(r.Context(), contextkey.UserID, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
