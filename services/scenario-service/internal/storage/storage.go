// работа с БД
package storage

import (
	"context"
	"database/sql"
	"scenario-service/internal/model"
)

type Storage struct {
	DB *sql.DB
}

func (s *Storage) GetScenarioByText(ctx context.Context, query string) (*model.Scenario, error) {
	var scenario model.Scenario
	err := s.DB.QueryRowContext(ctx, `
		SELECT id, title, location, price
		FROM scenarios
		WHERE title ILIKE '%' || $1 || '%' OR location ILIKE '%' || $1 || '%'
		LIMIT 1
	`, query).Scan(&scenario.ID, &scenario.Title, &scenario.Category, &scenario.Description)
	if err != nil {
		return nil, err
	}

	rows, err := s.DB.QueryContext(ctx, `
		SELECT id, scenario_id, text, place, time, price
		FROM steps
		WHERE scenario_id = $1
		ORDER BY time ASC
	`, scenario.ID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var step model.Step
		err := rows.Scan(&step.ID, &step.ScenarioID, &step.Text, &step.Place, &step.Time, &step.Price)
		if err != nil {
			return nil, err
		}
		scenario.Steps = append(scenario.Steps, step)
	}

	return &scenario, nil
}
