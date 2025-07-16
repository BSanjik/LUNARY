-- Удаляем индексы
DROP INDEX IF EXISTS idx_steps_order;
DROP INDEX IF EXISTS idx_steps_scenario_id;
DROP INDEX IF EXISTS idx_scenarios_category;
DROP INDEX IF EXISTS idx_scenarios_title;

-- Удаляем таблицы
DROP TABLE IF EXISTS steps;
DROP TABLE IF EXISTS scenarios;