package service

import (
	"auth-service/internal/model"
	"auth-service/internal/storage"
	"auth-service/internal/token"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Storage *storage.Storage
}

func (s *AuthService) Register(phone, password string) (*model.AuthResponse, error) {
	var isExist bool
	err := s.Storage.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE phone = $1)", phone).Scan(&isExist)
	if err != nil {
		return nil, err
	}
	if isExist {
		return nil, errors.New("пользователь с таким номером уже существует")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Создаем пользователя и получаем его данные
	var user model.User
	err = s.Storage.DB.QueryRow(
		"INSERT INTO users (phone, password) VALUES ($1, $2) RETURNING id, phone, email, created_at",
		phone, hashedPassword,
	).Scan(&user.ID, &user.Phone, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Генерируем токен
	tokenStr, err := token.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		Message: "Пользователь успешно зарегистрирован",
		Token:   tokenStr,
		User:    &user,
	}, nil
}

func (s *AuthService) Login(phone, password string) (*model.AuthResponse, error) {
	var user model.User
	var hashedPassword string

	err := s.Storage.DB.QueryRow(
		"SELECT id, phone, email, password, created_at FROM users WHERE phone = $1",
		phone,
	).Scan(&user.ID, &user.Phone, &user.Email, &hashedPassword, &user.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("неверный номер телефона или пароль")
	} else if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return nil, errors.New("неверный номер телефона или пароль")
	}

	// Генерируем токен
	tokenStr, err := token.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		Message: "Авторизация успешна",
		Token:   tokenStr,
		User:    &user,
	}, nil
}
