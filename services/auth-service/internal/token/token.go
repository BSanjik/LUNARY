// Работа с JWT
package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	Secret string
}

func New(secret string) *TokenService {
	return &TokenService{Secret: secret}
}

func (t *TokenService) Generate(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(t.Secret))
}

func (t *TokenService) Parse(tokenStr string) (int64, error) {
	token, err := jwt.Parse(tokenStr, func(tk *jwt.Token) (interface{}, error) {
		return []byte(t.Secret), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return int64(claims["user_id"].(float64)), nil
	}
	return 0, err
}
