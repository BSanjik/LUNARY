// HTTP обработчики (регистрация, логин)
package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/BSanjik/LUNARY/services/auth-service/internal/token"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DB *sql.DB
}

type registerRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type registerResponse struct {
	Token string `json:"token"`
}

var phoneRegex = regexp.MustCompile(`^\+?[0-9]{10,15}$`)

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		log.Println("[Register] JSON decode error:", err)
		return
	}

	req.Phone = strings.TrimSpace(req.Phone)
	req.Password = strings.TrimSpace(req.Password)
	req.Email = strings.TrimSpace(req.Email)

	//phone and pass validation
	if !phoneRegex.MatchString(req.Phone) {
		http.Error(w, "Неверный формат телефона", http.StatusBadRequest)
		return
	}
	if len(req.Password) < 12 {
		http.Error(w, "Пароль должен содержать 12 символов", http.StatusBadRequest)
		return
	}

	//isExist validation
	var isExist bool
	err := h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE phone=$1)", req.Phone).Scan(&isExist)
	if err != nil {
		http.Error(w, "DB Error ", http.StatusInternalServerError)
		log.Println("[Register] DB error:", err)
		return
	}
	if isExist {
		http.Error(w, "Пользователь с таким номером уже существует", http.StatusConflict)
		return
	}

	//hash pass
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	//insert DB
	var userID int64
	err = h.DB.QueryRow("INSERT INTO users (phone, password, email) VALUES ($1, $2, $3) RETURNING id", req.Phone, string(hashedPass), req.Email).Scan(&userID)
	if err != nil {
		http.Error(w, "Ошибка регистрации клиента", http.StatusInternalServerError)
		return
	}

	//JWT generation
	tokenStr, err := token.GenerateToken(userID)
	if err != nil {
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	//Response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(registerResponse{Token: tokenStr})

}
