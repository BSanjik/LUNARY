package service

import (
	"context"
	"scenario-service/internal/model"
)

// ScenarioServiceInterface определяет методы для работы со сценариями
type ScenarioServiceInterface interface {
	SearchScenarios(ctx context.Context, query string) (*model.Scenario, error)
	SubmitFeedback(ctx context.Context, feedback *model.FeedbackRequest) error
}
