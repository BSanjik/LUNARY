// бизнес-логика
package service

import (
	"context"
	"errors"
	"scenario-service/internal/model"
	"scenario-service/internal/storage"
)

type ScenarioService struct {
	Storage *storage.Storage
}

func (s *ScenarioService) GetScenario(ctx context.Context, query string) (*model.Scenario, error) {
	return s.Storage.GetScenarioByText(ctx, query)
}

func (s *ScenarioService) SubmitFeedback(ctx context.Context, feedback *model.FeedbackRequest) error {
	//
	if feedback.Rating < 1 || feedback.Rating > 5 {
		return errors.New("Рейтинг должен быть от 1 - 5")
	}

	//
	if feedback.ScenarioID <= 0 {
		return errors.New("Scenario_id пустой")
	}

	return s.Storage.SaveFeedback(ctx, feedback)
}
