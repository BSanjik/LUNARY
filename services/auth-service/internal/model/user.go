// Модели данных (пользователь)
package model

import "time"

type User struct {
	ID        int64     `json:"id"`
	Phone     string    `json:"phone"`
	Password  string    `json:"-"`     // не включаем в JSON ответы
	Email     *string   `json:"email"` // nullable поле
	CreatedAt time.Time `json:"created_at"`
}

// AuthResponse для возврата токена и пользователя
type AuthResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
	User    *User  `json:"user,omitempty"`
}
