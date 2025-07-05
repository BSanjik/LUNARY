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

type RegisterRequest struct {
	Phone    string `json:"phone" example:"+77001234567"`
	Password string `json:"password" example:"MySecurePassword123"`
	Email    string `json:"email" example:"user@example.com"`
}

type RegisterResponse struct {
	Token string `json:"token" example:"jwt.token.here"`
}

type LoginRequest struct {
	Phone    string `json:"phone" example:"+77001234567"`
	Password string `json:"password" example:"MySecurePassword123"`
}

type LoginResponse struct {
	Token string `json:"token" example:"jwt.token.here"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

var phoneRegex = regexp.MustCompile(`^\+?[0-9]{10,15}$`)

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(ErrorResponse{Message: message})
}

// Регистрация
// Register godoc
// @Summary Регистрация нового пользователя
// @Description Создает нового пользователя и возвращает JWT токен
// @Tags auth
// @Accept json
// @Produce json
// @Param input body RegisterRequest true "Данные пользователя"
// @Success 201 {object} RegisterResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("[Register] Ошибка декодирования JSON:", err)
		respondWithError(w, http.StatusBadRequest, "Неверный JSON")
		return
	}

	req.Phone = strings.TrimSpace(req.Phone)
	req.Password = strings.TrimSpace(req.Password)
	req.Email = strings.TrimSpace(req.Email)

	//phone and pass validation
	if !phoneRegex.MatchString(req.Phone) {
		respondWithError(w, http.StatusBadRequest, "Неверный формат телефона")
		return
	}
	if len(req.Password) < 12 {
		respondWithError(w, http.StatusBadRequest, "Пароль должен содержать минимум 12 символов")
		return
	}

	//isExist validation
	var isExist bool
	err := h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE phone=$1)", req.Phone).Scan(&isExist)
	if err != nil {
		log.Println("[Register] Ошибка запроса:", err)
		respondWithError(w, http.StatusInternalServerError, "Ошибка проверки пользователя в базе")
		return
	}
	if isExist {
		respondWithError(w, http.StatusConflict, "Пользователь с таким телефоном уже существует")
		return
	}

	//hash pass
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("[Register] bcrypt error:", err)
		respondWithError(w, http.StatusInternalServerError, "Ошибка шифрования пароля")
		return
	}

	//insert DB
	var userID int64
	err = h.DB.QueryRow("INSERT INTO users (phone, password, email) VALUES ($1, $2, $3) RETURNING id", req.Phone, string(hashedPass), req.Email).Scan(&userID)
	if err != nil {
		log.Println("[Register] Ошибка вставки:", err)
		respondWithError(w, http.StatusInternalServerError, "Ошибка сохранения пользователя")
		return
	}

	//JWT generation
	tokenStr, err := token.GenerateToken(userID)
	if err != nil {
		log.Println("[Register] Token error:", err)
		respondWithError(w, http.StatusInternalServerError, "Ошибка генерации токена")
		return
	}

	//Response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RegisterResponse{Token: tokenStr})

}

// Авторизация
// Login godoc
// @Summary Авторизация пользователя
// @Description Принимает телефон и пароль, возвращает JWT токен
// @Tags auth
// @Accept json
// @Produce json
// @Param input body LoginRequest true "Данные входа"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("[Login] Ошибка декодирования JSON:", err)
		respondWithError(w, http.StatusBadRequest, "Неверный JSON")
		return
	}

	req.Phone = strings.TrimSpace(req.Phone)
	req.Password = strings.TrimSpace(req.Password)

	if !phoneRegex.MatchString(req.Phone) {
		respondWithError(w, http.StatusBadRequest, "Неверный формат телефона")
		return
	}

	var userID int64
	var hashedPass string
	err := h.DB.QueryRow("SELECT id, password FROM users WHERE phone = $1", req.Phone).Scan(&userID, &hashedPass)
	if err == sql.ErrNoRows {
		respondWithError(w, http.StatusUnauthorized, "Неверный номер телефона или пароль")
		return
	} else if err != nil {
		log.Println("[Login] Ошибка запроса:", err)
		respondWithError(w, http.StatusInternalServerError, "Ошибка базы данных")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(req.Password))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Неверный номер телефона или пароль")
		return
	}

	tokenStr, err := token.GenerateToken(userID)
	if err != nil {
		log.Println("[Login] Token error:", err)
		respondWithError(w, http.StatusInternalServerError, "Ошибка генерации токена")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{Token: tokenStr})
}
