// бизнес-логика
package service

import (
	"context"
	"fmt"
	"scenario-service/internal/model"
	"scenario-service/internal/storage"
	"strings"

	"github.com/BSanjik/LUNARY/services/common/shared"
	"github.com/go-playground/validator/v10"
)

type ScenarioService struct {
	Storage   *storage.Storage
	Validator *validator.Validate
}

// DTO для валидации поискового запроса
type SearchDTO struct {
	Query string `json:"query" validate:"required,min=1,max=255"`
}

// DTO для валидации фидбека
type FeedbackDTO struct {
	ScenarioID int64  `json:"scenario_id" validate:"required,min=1"`
	Rating     int    `json:"rating" validate:"required,min=1,max=5"`
	Comment    string `json:"comment" validate:"max=1000"`
	UserLogin  string `json:"user_login" validate:"max=255"`
}

func (s *ScenarioService) SearchScenarios(ctx context.Context, query string) (*model.Scenario, error) {
	// Валидация входных данных
	dto := SearchDTO{Query: strings.TrimSpace(query)}
	if err := s.validateSearchData(dto); err != nil {
		shared.LogError(err, "Ошибка валидации поискового запроса", map[string]interface{}{
			"query": query,
		})
		return nil, err
	}

	shared.LogInfo("Поиск сценариев", map[string]interface{}{
		"query": dto.Query,
	})

	scenario, err := s.Storage.GetScenarioByText(ctx, dto.Query)
	if err != nil {
		shared.LogError(err, "Не удалось найти сценарий", map[string]interface{}{
			"query": dto.Query,
		})
		return nil, shared.NewNotFoundError("Сценарий не найден. Попробуйте изменить запрос")
	}

	shared.LogInfo("Сценарий найден", map[string]interface{}{
		"scenario_title": scenario.Title,
		"query":          dto.Query,
	})
	return scenario, nil
}

func (s *ScenarioService) SubmitFeedback(ctx context.Context, feedback *model.FeedbackRequest) error {
	// Валидация входных данных
	dto := FeedbackDTO{
		ScenarioID: feedback.ScenarioID,
		Rating:     feedback.Rating,
		Comment:    strings.TrimSpace(feedback.Comment),
		UserLogin:  strings.TrimSpace(feedback.UserLogin),
	}
	if err := s.validateFeedbackData(dto); err != nil {
		shared.LogError(err, "Ошибка валидации отзыва", map[string]interface{}{
			"scenario_id": feedback.ScenarioID,
			"rating":      feedback.Rating,
		})
		return err
	}

	shared.LogInfo("Сохранение отзыва", map[string]interface{}{
		"scenario_id": dto.ScenarioID,
		"rating":      dto.Rating,
	})

	// Обновляем исходную структуру с очищенными данными
	feedback.Comment = dto.Comment
	feedback.UserLogin = dto.UserLogin

	if err := s.Storage.SaveFeedback(ctx, feedback); err != nil {
		shared.LogError(err, "Не удалось сохранить отзыв в базе данных", map[string]interface{}{
			"scenario_id": feedback.ScenarioID,
		})
		return shared.NewInternalError("Не удалось сохранить отзыв")
	}

	shared.LogInfo("Отзыв успешно сохранён", map[string]interface{}{
		"scenario_id": feedback.ScenarioID,
	})
	return nil
}

func (s *ScenarioService) validateSearchData(dto SearchDTO) error {
	if err := s.Validator.Struct(dto); err != nil {
		return s.formatValidationError(err)
	}
	return nil
}

func (s *ScenarioService) validateFeedbackData(dto FeedbackDTO) error {
	if err := s.Validator.Struct(dto); err != nil {
		return s.formatValidationError(err)
	}
	return nil
}

func (s *ScenarioService) formatValidationError(err error) error {
	var messages []string
	for _, err := range err.(validator.ValidationErrors) {
		message := s.getValidationMessage(err)
		messages = append(messages, message)
	}
	return shared.NewValidationError(strings.Join(messages, "; "))
}

func (s *ScenarioService) getValidationMessage(err validator.FieldError) string {
	field := err.Field()
	tag := err.Tag()

	switch field {
	case "Query":
		switch tag {
		case "required":
			return "Поисковый запрос обязателен для заполнения"
		case "min":
			return "Поисковый запрос не может быть пустым"
		case "max":
			return "Поисковый запрос слишком длинный (максимум 255 символов)"
		}
	case "ScenarioID":
		switch tag {
		case "required":
			return "ID сценария обязателен"
		case "min":
			return "Некорректный ID сценария"
		}
	case "Rating":
		switch tag {
		case "required":
			return "Рейтинг обязателен"
		case "min":
			return "Рейтинг должен быть от 1 до 5"
		case "max":
			return "Рейтинг должен быть от 1 до 5"
		}
	case "Comment":
		switch tag {
		case "max":
			return "Комментарий слишком длинный (максимум 1000 символов)"
		}
	case "UserLogin":
		switch tag {
		case "max":
			return "Логин пользователя слишком длинный (максимум 255 символов)"
		}
	}

	return fmt.Sprintf("Некорректное значение поля %s", field)
}
