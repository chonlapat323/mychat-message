package utils

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims คือ payload ของ token ที่เรากำหนดเอง
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken สร้าง JWT token สำหรับผู้ใช้คนหนึ่ง
func GenerateTokens(userID, email, role string) (accessToken string, refreshToken string, err error) {
	secret := os.Getenv("JWT_SECRET")

	// Access Token: อายุสั้น
	accessExpire := time.Now().Add(15 * time.Minute)
	accessClaims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpire),
		},
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = at.SignedString([]byte(secret))
	if err != nil {
		return
	}

	// Refresh Token: อายุยาว
	refreshExpire := time.Now().Add(7 * 24 * time.Hour)
	refreshClaims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpire),
		},
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = rt.SignedString([]byte(secret))
	return
}

// ValidateToken ถอดรหัสและตรวจสอบ JWT token
func ValidateToken(tokenStr string) (*Claims, error) {
	log.Println("🧪 ValidateToken called")

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Println("❌ JWT_SECRET is empty")
		return nil, errors.New("missing JWT_SECRET")
	}

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		log.Println("❌ Token parsing failed:", err)
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		log.Println("❌ Token claims type assertion failed")
		return nil, errors.New("invalid token claims")
	}
	if !token.Valid {
		log.Println("❌ Token is not valid")
		return nil, errors.New("invalid token")
	}

	log.Println("✅ Token is valid for user:", claims.UserID)
	return claims, nil
}
