package service

import (
	"auth-service/internal/model"
	"auth-service/internal/storage"
	"auth-service/internal/token"
	"database/sql"
	"fmt"
	"strings"

	"github.com/BSanjik/LUNARY/services/common/shared"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Storage   *storage.Storage
	Validator *validator.Validate
}

// DTO для валидации регистрации
type RegisterDTO struct {
	Phone    string `validate:"required,e164"`
	Password string `validate:"required,min=6"`
}

// DTO для валидации входа
type LoginDTO struct {
	Phone    string `validate:"required,e164"`
	Password string `validate:"required"`
}

// Валидация данных регистрации
func (s *AuthService) validateRegistrationData(phone, password string) error {
	dto := RegisterDTO{
		Phone:    strings.TrimSpace(phone),
		Password: strings.TrimSpace(password),
	}

	if err := s.Validator.Struct(dto); err != nil {
		return s.formatValidationError(err)
	}

	return nil
}

// Валидация данных входа
func (s *AuthService) validateLoginData(phone, password string) error {
	dto := LoginDTO{
		Phone:    strings.TrimSpace(phone),
		Password: strings.TrimSpace(password),
	}

	if err := s.Validator.Struct(dto); err != nil {
		return s.formatValidationError(err)
	}

	return nil
}

func (s *AuthService) formatValidationError(err error) error {
	var details []string
	for _, validErr := range err.(validator.ValidationErrors) {
		details = append(details, fmt.Sprintf("Поле '%s': %s",
			validErr.Field(), s.getValidationMessage(validErr)))
	}
	return shared.NewValidationError("Ошибка валидации: " + strings.Join(details, "; "))
}

func (s *AuthService) getValidationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "обязательно для заполнения"
	case "email":
		return "должно быть валидным email"
	case "min":
		return "должно содержать минимум 6 символов"
	case "e164":
		return "должно быть в международном формате (+77001234567)"
	default:
		return "неверный формат"
	}
}

func (s *AuthService) Register(phone, password string) (*model.AuthResponse, error) {
	// Валидация входных данных
	if err := s.validateRegistrationData(phone, password); err != nil {
		shared.LogAuthEvent("register", phone, 0, false, "validation_failed")
		return nil, err
	}

	// Применяем trim (валидация уже проверила, но нормализуем данные)
	phone = strings.TrimSpace(phone)
	password = strings.TrimSpace(password)

	// Проверяем существование пользователя
	var isExist bool
	err := s.Storage.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE phone = $1)", phone).Scan(&isExist)
	if err != nil {
		shared.LogError(err, "Database error during user existence check", map[string]interface{}{
			"phone":     phone,
			"operation": "user_exists_check",
		})
		return nil, shared.NewInternalError("Ошибка проверки пользователя в базе данных")
	}

	if isExist {
		shared.LogAuthEvent("register", phone, 0, false, "user_already_exists")
		return nil, shared.NewConflictError("Пользователь с таким номером уже существует")
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		shared.LogError(err, "Password hashing failed", map[string]interface{}{
			"phone": phone,
		})
		return nil, shared.NewInternalError("Ошибка хеширования пароля")
	}

	// Создаем пользователя
	var user model.User
	err = s.Storage.DB.QueryRow(
		"INSERT INTO users (phone, password) VALUES ($1, $2) RETURNING id, phone, email, created_at",
		phone, hashedPassword,
	).Scan(&user.ID, &user.Phone, &user.CreatedAt)

	if err != nil {
		shared.LogError(err, "User creation failed", map[string]interface{}{
			"phone":     phone,
			"operation": "create_user",
		})
		return nil, shared.NewInternalError("Ошибка создания пользователя")
	}

	// Генерируем токен
	tokenStr, err := token.GenerateToken(user.ID)
	if err != nil {
		shared.LogError(err, "Token generation failed", map[string]interface{}{
			"user_id": user.ID,
			"phone":   phone,
		})
		return nil, shared.NewInternalError("Ошибка генерации токена")
	}

	// Успешная регистрация
	shared.LogAuthEvent("register", phone, user.ID, true, "")
	shared.LogInfo("User registered successfully", map[string]interface{}{
		"user_id": user.ID,
		"phone":   phone,
	})

	return &model.AuthResponse{
		Message: "Пользователь успешно зарегистрирован",
		Token:   tokenStr,
		User:    &user,
	}, nil
}

func (s *AuthService) Login(phone, password string) (*model.AuthResponse, error) {
	// Валидация входных данных
	if err := s.validateLoginData(phone, password); err != nil {
		shared.LogAuthEvent("login", phone, 0, false, "validation_failed")
		return nil, err
	}

	// Применяем trim (валидация уже проверила, но нормализуем данные)
	phone = strings.TrimSpace(phone)
	password = strings.TrimSpace(password)

	var user model.User
	var hashedPassword string

	err := s.Storage.DB.QueryRow(
		"SELECT id, phone, email, password, created_at FROM users WHERE phone = $1",
		phone,
	).Scan(&user.ID, &user.Phone, &hashedPassword, &user.CreatedAt)

	if err == sql.ErrNoRows {
		shared.LogAuthEvent("login", phone, 0, false, "user_not_found")
		return nil, shared.NewAuthError("Неверный номер телефона или пароль")
	} else if err != nil {
		shared.LogError(err, "Database error during user lookup", map[string]interface{}{
			"phone":     phone,
			"operation": "user_lookup",
		})
		return nil, shared.NewInternalError("Ошибка поиска пользователя")
	}

	// Проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		shared.LogAuthEvent("login", phone, user.ID, false, "invalid_password")
		return nil, shared.NewAuthError("Неверный номер телефона или пароль")
	}

	// Генерируем токен
	tokenStr, err := token.GenerateToken(user.ID)
	if err != nil {
		shared.LogError(err, "Token generation failed during login", map[string]interface{}{
			"user_id": user.ID,
			"phone":   phone,
		})
		return nil, shared.NewInternalError("Ошибка генерации токена")
	}

	// Успешный вход
	shared.LogAuthEvent("login", phone, user.ID, true, "")
	shared.LogInfo("User logged in successfully", map[string]interface{}{
		"user_id": user.ID,
		"phone":   phone,
	})

	return &model.AuthResponse{
		Message: "Авторизация успешна",
		Token:   tokenStr,
		User:    &user,
	}, nil
}
