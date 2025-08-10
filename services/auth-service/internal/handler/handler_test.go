package handler

import (
	"auth-service/internal/model"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(phone, password string) (*model.AuthResponse, error) {
	args := m.Called(phone, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.AuthResponse), args.Error(1)
}

func (m *MockAuthService) Login(phone, password string) (*model.AuthResponse, error) {
	args := m.Called(phone, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.AuthResponse), args.Error(1)
}

func TestHandler_Registration_Success(t *testing.T) {
	mockService := new(MockAuthService)
	handler := &Handler{Service: mockService}

	// Настраиваем mock
	expectedResponse := &model.AuthResponse{
		Message: "Пользователь успешно зарегистрирован",
		Token:   "test-token",
		User: &model.User{
			ID:    1,
			Phone: "+77001234567",
		},
	}
	mockService.On("Register", "+77001234567", "password123").Return(expectedResponse, nil)

	// Подготавливаем запрос
	reqBody := RegisterRequest{
		Phone:    "+77001234567",
		Password: "password123",
		Email:    "test@example.com",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/registration", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// Выполняем запрос
	handler.Registration(rr, req)

	// Проверяем результат
	assert.Equal(t, http.StatusCreated, rr.Code)

	var response model.AuthResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "test-token", response.Token)

	mockService.AssertExpectations(t)
}

func TestHandler_Registration_InvalidPhone(t *testing.T) {
	mockService := new(MockAuthService)
	handler := &Handler{Service: mockService}

	// Подготавливаем запрос с неверным телефоном
	reqBody := RegisterRequest{
		Phone:    "invalid-phone",
		Password: "password123",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/registration", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// Выполняем запрос
	handler.Registration(rr, req)

	// Проверяем результат
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Неверный формат телефона")
}
