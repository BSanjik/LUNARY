-- Создание таблицы сценариев
CREATE TABLE IF NOT EXISTS scenarios (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    category TEXT,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы шагов
CREATE TABLE IF NOT EXISTS steps (
    id SERIAL PRIMARY KEY,
    scenario_id INTEGER REFERENCES scenarios(id) ON DELETE CASCADE,
    text TEXT NOT NULL,
    place TEXT,
    time TEXT,
    price DECIMAL(10,2) DEFAULT 0,
    step_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Добавляем индексы для оптимизации
CREATE INDEX IF NOT EXISTS idx_scenarios_title ON scenarios(title);
CREATE INDEX IF NOT EXISTS idx_scenarios_category ON scenarios(category);
CREATE INDEX IF NOT EXISTS idx_steps_scenario_id ON steps(scenario_id);
CREATE INDEX IF NOT EXISTS idx_steps_order ON steps(scenario_id, step_order);