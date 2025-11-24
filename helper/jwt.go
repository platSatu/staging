package helper

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var accessTokenTTL = time.Minute * 15 // 15 menit

// ======================== GENERATE ACCESS TOKEN ========================
func GenerateAccessToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(accessTokenTTL).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "b6c0f8a23e9f4d7e8b3a2d1f0c9e4a7d9b0c1e2f3d4a5b6c7d8e9f0a1b2c3d4e"
	}
	if len(secret) < 32 {
		return "", errors.New("JWT_SECRET harus minimal 256-bit / 32 karakter")
	}

	return token.SignedString([]byte(secret))
}

// ======================== VALIDATE TOKEN ========================
func ValidateAccessToken(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "b6c0f8a23e9f4d7e8b3a2d1f0c9e4a7d9b0c1e2f3d4a5b6c7d8e9f0a1b2c3d4e"
	}
	if len(secret) < 32 {
		return nil, errors.New("JWT_SECRET harus minimal 256-bit / 32 karakter")
	}

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
}

// ======================== GET USER ID FROM TOKEN ========================
func GetUserIDFromToken(tokenString string) (string, error) {
	token, err := ValidateAccessToken(tokenString)
	if err != nil || !token.Valid {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		return "", errors.New("userID tidak ditemukan di token")
	}

	return userID, nil
}
