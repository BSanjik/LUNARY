// HTTP обработчики (регистрация, логин)
package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/BSanjik/LUNARY/services/auth-service/internal/hash"
	"github.com/BSanjik/LUNARY/services/auth-service/internal/logger"
	"github.com/BSanjik/LUNARY/services/auth-service/internal/token"
	"github.com/BSanjik/LUNARY/services/auth-service/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DB           *sql.DB
	TokenService *token.TokenService
	HashService  *hash.HashService
}

func NewAuthHandler(db *sql.DB, tokenService *token.TokenService) *AuthHandler {
	return &AuthHandler{
		DB:           db,
		TokenService: tokenService,
		HashService:  hash.New(),
	}
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
		logger.Log.Errorw("Ошибка декодирования JSON", "error", err)
		utils.JSON(w, http.StatusBadRequest, nil, "Неверный JSON")
		return
	}

	req.Phone = strings.TrimSpace(req.Phone)
	req.Password = strings.TrimSpace(req.Password)
	req.Email = strings.TrimSpace(req.Email)

	//phone and pass validation
	if !phoneRegex.MatchString(req.Phone) {
		utils.JSON(w, http.StatusBadRequest, nil, "Неверный формат телефона")
		return
	}
	if len(req.Password) < 12 {
		utils.JSON(w, http.StatusBadRequest, nil, "Пароль должен содержать минимум 12 символов")
		return
	}

	//isExist validation
	var isExist bool
	err := h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE phone=$1)", req.Phone).Scan(&isExist)
	if err != nil {
		logger.Log.Errorw("Ошибка запроса в БД (exists)", "error", err)
		utils.JSON(w, http.StatusInternalServerError, nil, "Ошибка проверки пользователя")
		return
	}
	if isExist {
		utils.JSON(w, http.StatusConflict, nil, "Пользователь с таким телефоном уже существует")
		return
	}

	//hash pass
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Errorw("Ошибка хеширования пароля", "error", err)
		utils.JSON(w, http.StatusInternalServerError, nil, "Ошибка хеширования")
		return
	}

	//insert DB
	var userID int64
	err = h.DB.QueryRow("INSERT INTO users (phone, password, email) VALUES ($1, $2, $3) RETURNING id", req.Phone, string(hashedPass), req.Email).Scan(&userID)
	if err != nil {
		logger.Log.Errorw("Ошибка вставки в БД", "error", err)
		utils.JSON(w, http.StatusInternalServerError, nil, "Ошибка сохранения пользователя")
		return
	}

	//JWT generation
	tokenStr, err := h.TokenService.Generate(userID)
	if err != nil {
		logger.Log.Errorw("Ошибка генерации токена", "error", err)
		utils.JSON(w, http.StatusInternalServerError, nil, "Ошибка генерации токена")
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
		logger.Log.Errorw("Ошибка декодирования JSON", "error", err)
		utils.JSON(w, http.StatusBadRequest, nil, "Неверный JSON")
		return
	}

	req.Phone = strings.TrimSpace(req.Phone)
	req.Password = strings.TrimSpace(req.Password)

	if !phoneRegex.MatchString(req.Phone) {
		utils.JSON(w, http.StatusBadRequest, nil, "Неверный формат телефона")
		return
	}

	var userID int64
	var hashedPass string
	err := h.DB.QueryRow("SELECT id, password FROM users WHERE phone = $1", req.Phone).Scan(&userID, &hashedPass)
	if err == sql.ErrNoRows {
		utils.JSON(w, http.StatusUnauthorized, nil, "Неверный номер телефона или пароль")
		return
	} else if err != nil {
		logger.Log.Errorw("Ошибка запроса в БД", "error", err)
		utils.JSON(w, http.StatusInternalServerError, nil, "Ошибка базы данных")
		return
	}

	if !h.HashService.CheckPasswordHash(req.Password, hashedPass) {
		utils.JSON(w, http.StatusUnauthorized, nil, "Неверный номер телефона или пароль")
		return
	}

	tokenStr, err := h.TokenService.Generate(userID)
	if err != nil {
		logger.Log.Errorw("Ошибка генерации токена", "error", err)
		utils.JSON(w, http.StatusInternalServerError, nil, "Ошибка генерации токена")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{Token: tokenStr})
}
