package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"scenario-service/internal/model"
	"testing"

	"github.com/BSanjik/LUNARY/services/common/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockScenarioService реализует ScenarioServiceInterface для тестирования
type MockScenarioService struct {
	mock.Mock
}

func (m *MockScenarioService) SearchScenarios(ctx context.Context, query string) (*model.Scenario, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Scenario), args.Error(1)
}

func (m *MockScenarioService) SubmitFeedback(ctx context.Context, feedback *model.FeedbackRequest) error {
	args := m.Called(ctx, feedback)
	return args.Error(0)
}

func TestHandler_GetScenario(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(*MockScenarioService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Successful search",
			requestBody: scenarioRequest{
				Query: "семейный ужин",
			},
			setupMock: func(mockService *MockScenarioService) {
				scenario := &model.Scenario{
					ID:          1,
					Title:       "Семейный ужин дома",
					Category:    "семейный",
					Description: "Уютный ужин в кругу семьи",
				}
				mockService.On("SearchScenarios", mock.Anything, "семейный ужин").Return(scenario, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"title":"Семейный ужин дома","category":"семейный","description":"Уютный ужин в кругу семьи","steps":null}`,
		},
		{
			name: "Scenario not found",
			requestBody: scenarioRequest{
				Query: "несуществующий сценарий",
			},
			setupMock: func(mockService *MockScenarioService) {
				mockService.On("SearchScenarios", mock.Anything, "несуществующий сценарий").Return(nil, shared.NewNotFoundError("Сценарий не найден"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"Сценарий не найден","code":"NOT_FOUND"}`,
		},
		{
			name:           "Invalid JSON format",
			requestBody:    "invalid json",
			setupMock:      func(mockService *MockScenarioService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Некорректный формат запроса","code":"VALIDATION_ERROR"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockService := new(MockScenarioService)
			tt.setupMock(mockService)

			handler := &Handler{Service: mockService}

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/search", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			// Act
			handler.GetScenario(rr, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.JSONEq(t, tt.expectedBody, rr.Body.String())
			mockService.AssertExpectations(t)
		})
	}
}

func TestHandler_SubmitFeedback(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(*MockScenarioService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Successful feedback submission",
			requestBody: model.FeedbackRequest{
				ScenarioID: 1,
				Rating:     5,
				Comment:    "Отличный сценарий!",
				UserLogin:  "testuser",
			},
			setupMock: func(mockService *MockScenarioService) {
				mockService.On("SubmitFeedback", mock.Anything, mock.AnythingOfType("*model.FeedbackRequest")).Return(nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"status":"success","message":"Отзыв успешно сохранён"}`,
		},
		{
			name: "Invalid rating",
			requestBody: model.FeedbackRequest{
				ScenarioID: 1,
				Rating:     6,
				Comment:    "Плохой рейтинг",
			},
			setupMock: func(mockService *MockScenarioService) {
				mockService.On("SubmitFeedback", mock.Anything, mock.AnythingOfType("*model.FeedbackRequest")).Return(shared.NewValidationError("Рейтинг должен быть от 1 до 5"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Рейтинг должен быть от 1 до 5","code":"VALIDATION_ERROR"}`,
		},
		{
			name:           "Invalid JSON format",
			requestBody:    "invalid json",
			setupMock:      func(mockService *MockScenarioService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Некорректный формат данных","code":"VALIDATION_ERROR"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockService := new(MockScenarioService)
			tt.setupMock(mockService)

			handler := &Handler{Service: mockService}

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/feedback", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			// Act
			handler.SubmitFeedback(rr, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.JSONEq(t, tt.expectedBody, rr.Body.String())
			mockService.AssertExpectations(t)
		})
	}
}
