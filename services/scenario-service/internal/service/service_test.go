package service

import (
	"context"
	"database/sql"
	"scenario-service/internal/model"
	"scenario-service/internal/storage"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestScenarioService_SearchScenarios(t *testing.T) {
	// Создаём mock БД
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	store := &storage.Storage{DB: db}
	validator := validator.New()
	service := &ScenarioService{
		Storage:   store,
		Validator: validator,
	}

	tests := []struct {
		name      string
		query     string
		setupMock func()
		wantErr   bool
		errMsg    string
	}{
		{
			name:  "Successful search",
			query: "семейный ужин",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "category", "description"}).
					AddRow(1, "Семейный ужин дома", "семейный", "Уютный ужин в кругу семьи")
				mock.ExpectQuery(`SELECT .+ FROM scenarios`).
					WithArgs("семейный ужин").
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name:  "Empty query validation error",
			query: "",
			setupMock: func() {
				// No DB expectations since validation should fail first
			},
			wantErr: true,
			errMsg:  "Поисковый запрос обязателен для заполнения",
		},
		{
			name:  "Query too long validation error",
			query: string(make([]byte, 300)), // More than 255 characters
			setupMock: func() {
				// No DB expectations since validation should fail first
			},
			wantErr: true,
			errMsg:  "Поисковый запрос слишком длинный (максимум 255 символов)",
		},
		{
			name:  "Scenario not found",
			query: "несуществующий сценарий",
			setupMock: func() {
				mock.ExpectQuery(`SELECT .+ FROM scenarios`).
					WithArgs("несуществующий сценарий").
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
			errMsg:  "Сценарий не найден. Попробуйте изменить запрос",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			result, err := service.SearchScenarios(context.Background(), tt.query)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, int64(1), result.ID)
				assert.Equal(t, "Семейный ужин дома", result.Title)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestScenarioService_SubmitFeedback(t *testing.T) {
	// Создаём mock БД
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	store := &storage.Storage{DB: db}
	validator := validator.New()
	service := &ScenarioService{
		Storage:   store,
		Validator: validator,
	}

	tests := []struct {
		name      string
		feedback  *model.FeedbackRequest
		setupMock func()
		wantErr   bool
		errMsg    string
	}{
		{
			name: "Successful feedback submission",
			feedback: &model.FeedbackRequest{
				ScenarioID: 1,
				Rating:     5,
				Comment:    "Отличный сценарий!",
				UserLogin:  "testuser",
			},
			setupMock: func() {
				mock.ExpectExec(`INSERT INTO feedback`).
					WithArgs(1, 5, "Отличный сценарий!", "testuser").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "Invalid rating - too low",
			feedback: &model.FeedbackRequest{
				ScenarioID: 1,
				Rating:     0,
				Comment:    "Плохой рейтинг",
			},
			setupMock: func() {
				// No DB expectations since validation should fail first
			},
			wantErr: true,
			errMsg:  "Рейтинг должен быть от 1 до 5",
		},
		{
			name: "Invalid rating - too high",
			feedback: &model.FeedbackRequest{
				ScenarioID: 1,
				Rating:     6,
				Comment:    "Слишком высокий рейтинг",
			},
			setupMock: func() {
				// No DB expectations since validation should fail first
			},
			wantErr: true,
			errMsg:  "Рейтинг должен быть от 1 до 5",
		},
		{
			name: "Invalid scenario ID",
			feedback: &model.FeedbackRequest{
				ScenarioID: 0,
				Rating:     5,
				Comment:    "Некорректный ID",
			},
			setupMock: func() {
				// No DB expectations since validation should fail first
			},
			wantErr: true,
			errMsg:  "Некорректный ID сценария",
		},
		{
			name: "Comment too long",
			feedback: &model.FeedbackRequest{
				ScenarioID: 1,
				Rating:     5,
				Comment:    string(make([]byte, 1100)), // More than 1000 characters
			},
			setupMock: func() {
				// No DB expectations since validation should fail first
			},
			wantErr: true,
			errMsg:  "Комментарий слишком длинный (максимум 1000 символов)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := service.SubmitFeedback(context.Background(), tt.feedback)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
