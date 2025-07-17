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
		SELECT id, title, category, description
		FROM scenarios
		WHERE title ILIKE '%' || $1 || '%' OR description ILIKE '%' || $1 || '%' OR category ILIKE '%' || $1 || '%'
		LIMIT 1
	`, query).Scan(&scenario.ID, &scenario.Title, &scenario.Category, &scenario.Description)
	if err != nil {
		return nil, err
	}

	rows, err := s.DB.QueryContext(ctx, `
		SELECT id, scenario_id, text, place, time, price
		FROM steps
		WHERE scenario_id = $1
		ORDER BY step_order ASC
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

func (s *Storage) SaveFeedback(ctx context.Context, feedback *model.FeedbackRequest) error {
	_, err := s.DB.ExecContext(ctx, `
		INSERT INTO feedback (scenario_id, rating, comment, user_email, created_at)
		VALUES ($1, $2, $3, $4, NOW())
	`, feedback.ScenarioID, feedback.Rating, feedback.Comment, feedback.UserLogin)

	return err
}
