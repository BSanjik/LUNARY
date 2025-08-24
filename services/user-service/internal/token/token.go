// JWT токены и авторизация
package token

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int64  `json:"user_id"`
	Phone  string `json:"phone"`
	jwt.RegisteredClaims
}

var jwtSecret string

// Инициализация секретного ключа
func Init(secret string) {
	jwtSecret = secret
}

// Парсинг и валидация JWT токена
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// Контекстные ключи
type contextKey string

const (
	userIDKey contextKey = "user_id"
	phoneKey  contextKey = "phone"
)

// Добавить user_id в контекст
func ContextWithUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// Добавить phone в контекст
func ContextWithPhone(ctx context.Context, phone string) context.Context {
	return context.WithValue(ctx, phoneKey, phone)
}

// Получить user_id из контекста
func UserIDFromContext(ctx context.Context) (int64, error) {
	userID, ok := ctx.Value(userIDKey).(int64)
	if !ok {
		return 0, fmt.Errorf("user_id not found in context")
	}
	return userID, nil
}

// Получить phone из контекста
func PhoneFromContext(ctx context.Context) (string, error) {
	phone, ok := ctx.Value(phoneKey).(string)
	if !ok {
		return "", fmt.Errorf("phone not found in context")
	}
	return phone, nil
}
