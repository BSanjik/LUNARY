// бизнес-логика
package service

import (
	"context"
	"scenario-service/internal/model"
	"scenario-service/internal/storage"
)

type ScenarioService struct {
	Storage *storage.Storage
}

func (s *ScenarioService) GetScenario(ctx context.Context, query string) (*model.Scenario, error) {
	return s.Storage.GetScenarioByText(ctx, query)
}
