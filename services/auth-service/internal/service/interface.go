package service

import "auth-service/internal/model"

type AuthServiceInterface interface {
	Register(phone, password string) (*model.AuthResponse, error)
	Login(phone, password string) (*model.AuthResponse, error)
}
