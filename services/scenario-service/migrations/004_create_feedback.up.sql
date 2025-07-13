-- Создание таблицы для фидбека
CREATE TABLE IF NOT EXISTS feedback (
    id SERIAL PRIMARY KEY,
    scenario_id INTEGER REFERENCES scenarios(id) ON DELETE CASCADE,
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    user_email VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индекс для быстрого поиска фидбека по сценарию
CREATE INDEX IF NOT EXISTS idx_feedback_scenario_id ON feedback(scenario_id);
CREATE INDEX IF NOT EXISTS idx_feedback_rating ON feedback(rating);