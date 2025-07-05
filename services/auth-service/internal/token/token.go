// Работа с JWT
package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JWTkey []byte

func Init(secret string) {
	JWTkey = []byte(secret)
}

func GenerateToken(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTkey)
}

func ValidateToken(tokenStr string) (int64, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return JWTkey, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return int64(claims["user_id"].(float64)), nil
	}
	return 0, err
}
