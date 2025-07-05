// HTTP обработчики (регистрация, логин)
package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/BSanjik/LUNARY/services/auth-service/internal/model"
	"github.com/BSanjik/LUNARY/services/auth-service/internal/token"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DB *sql.DB
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var u model.User
	json.NewDecoder(r.Body).Decode(&u)
	hashed, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	err := h.DB.QueryRow("INSERT INTO users(phone, password, email) VALUES($1, $2) RETURNING id", u.Phone, string(hashed), u.Email).Scan(&u.ID)
	if err != nil {
		http.Error(w, "Email already exists", http.StatusBadRequest)
		return
	}
	tokenStr, _ := token.GenerateToken(u.ID)
	w.Write([]byte(tokenStr))
}
