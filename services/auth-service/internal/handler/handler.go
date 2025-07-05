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

type loginRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

var phoneRegex = regexp.MustCompile(`^\+?[0-9]{10,15}$`)

// Регистрация
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		log.Println("[Register] Ошибка декодирования JSON:", err)
		return
	}

	req.Phone = strings.TrimSpace(req.Phone)
	req.Password = strings.TrimSpace(req.Password)
	req.Email = strings.TrimSpace(req.Email)

	//phone and pass validation
	if !phoneRegex.MatchString(req.Phone) {
		http.Error(w, "Неверный формат телефона. Используйте только цифры, от 10 до 15 символов", http.StatusBadRequest)
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
		http.Error(w, "Ошибка при проверке телефона в базе ", http.StatusInternalServerError)
		log.Println("[Register] Ошибка запроса:", err)
		return
	}
	if isExist {
		http.Error(w, "Пользователь с таким телефоном уже зарегистрирован", http.StatusConflict)
		return
	}

	//hash pass
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Ошибка при хешировании пароля", http.StatusInternalServerError)
		log.Println("[Register] bcrypt error:", err)
		return
	}

	//insert DB
	var userID int64
	err = h.DB.QueryRow("INSERT INTO users (phone, password, email) VALUES ($1, $2, $3) RETURNING id", req.Phone, string(hashedPass), req.Email).Scan(&userID)
	if err != nil {
		http.Error(w, "Ошибка при сохранении пользователя", http.StatusInternalServerError)
		log.Println("[Register] Ошибка вставки:", err)
		return
	}

	//JWT generation
	tokenStr, err := token.GenerateToken(userID)
	if err != nil {
		http.Error(w, "Ошибка генерации токена", http.StatusInternalServerError)
		log.Println("[Register] Token error:", err)
		return
	}

	//Response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(registerResponse{Token: tokenStr})

}

// Авторизация
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		log.Println("[Login] Ошибка декодирования JSON:", err)
		return
	}

	req.Phone = strings.TrimSpace(req.Phone)
	req.Password = strings.TrimSpace(req.Password)

	if !phoneRegex.MatchString(req.Phone) {
		http.Error(w, "Неверный формат телефона", http.StatusBadRequest)
		return
	}

	var userID int64
	var hashedPass string
	err := h.DB.QueryRow("SELECT id, password FROM users WHERE phone = $1", req.Phone).Scan(&userID, &hashedPass)
	if err == sql.ErrNoRows {
		http.Error(w, "Неверный номер телефона или пароль", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Ошибка при получении пользователя", http.StatusInternalServerError)
		log.Println("[Login] Ошибка запроса:", err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(req.Password))
	if err != nil {
		http.Error(w, "Неверный номер телефона или пароль", http.StatusUnauthorized)
		return
	}

	tokenStr, err := token.GenerateToken(userID)
	if err != nil {
		http.Error(w, "Ошибка генерации токена", http.StatusInternalServerError)
		log.Println("[Login] Token error:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loginResponse{Token: tokenStr})
}
