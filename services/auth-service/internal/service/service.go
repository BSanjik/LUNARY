package service

import (
	"auth-service/internal/storage"
	"auth-service/internal/token"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Storage *storage.Storage
}

func (s *AuthService) Register(phone, password string) (string, error) {
	var isExist bool
	err := s.Storage.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE phone = $1)", phone).Scan(&isExist)
	if err != nil {
		return "", err
	}
	if isExist {
		return "", errors.New("пользователь с таким номером уже существует")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	var userID int64
	err = s.Storage.DB.QueryRow("INSERT INTO users (phone, password) VALUES ($1, $2) RETURNING id", phone, hashedPassword).Scan(&userID)
	if err != nil {
		return "", err
	}

	return token.GenerateToken(userID)
}

func (s *AuthService) Login(phone, password string) (string, error) {
	var userID int64
	var hashedPassword string
	err := s.Storage.DB.QueryRow("SELECT id, password FROM users WHERE phone = $1", phone).Scan(&userID, &hashedPassword)
	if err == sql.ErrNoRows {
		return "", errors.New("неверный номер телефона или пароль")
	} else if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return "", errors.New("неверный номер телефона или пароль")
	}

	return token.GenerateToken(userID)
}
