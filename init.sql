-- Единая схема базы данных LUNARY для MVP
-- Объединяет таблицы для всех микросервисов

-- ========================================
-- AUTH SERVICE TABLES
-- ========================================

-- Таблица пользователей
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    phone VARCHAR(20) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ========================================
-- SCENARIO SERVICE TABLES
-- ========================================

-- Таблица сценариев
CREATE TABLE IF NOT EXISTS scenarios (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    category TEXT,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица шагов сценариев
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

-- Таблица фидбека
CREATE TABLE IF NOT EXISTS feedback (
    id SERIAL PRIMARY KEY,
    scenario_id INTEGER REFERENCES scenarios(id) ON DELETE CASCADE,
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    user_email VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ========================================
-- INDEXES FOR OPTIMIZATION
-- ========================================

-- Индексы для таблицы пользователей
CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Индексы для таблицы сценариев
CREATE INDEX IF NOT EXISTS idx_scenarios_title ON scenarios(title);
CREATE INDEX IF NOT EXISTS idx_scenarios_category ON scenarios(category);

-- Индексы для таблицы шагов
CREATE INDEX IF NOT EXISTS idx_steps_scenario_id ON steps(scenario_id);
CREATE INDEX IF NOT EXISTS idx_steps_order ON steps(scenario_id, step_order);

-- Индексы для таблицы фидбека
CREATE INDEX IF NOT EXISTS idx_feedback_scenario_id ON feedback(scenario_id);
CREATE INDEX IF NOT EXISTS idx_feedback_rating ON feedback(rating);

-- ========================================
-- SAMPLE DATA FOR TESTING
-- ========================================

-- Тестовые сценарии
INSERT INTO scenarios (id, title, category, description) VALUES 
(1, 'Семейный выходной в Алматы', 'семья', 'Идеальный день для семьи с детьми: парк, кафе, развлечения'),
(2, 'Корпоративный тимбилдинг', 'корпоратив', 'Активный корпоративный день для команды 15-20 человек')
ON CONFLICT (id) DO NOTHING;

-- Шаги для семейного выходного
INSERT INTO steps (scenario_id, text, place, time, price, step_order) VALUES
(1, 'Завтрак в семейном кафе', 'Кафе "Семейное гнездышко", ул. Абая 150', '09:00-10:00', 8000.00, 1),
(1, 'Прогулка в Центральном парке', 'Центральный парк культуры и отдыха', '10:30-12:00', 0.00, 2),
(1, 'Обед в семейном ресторане', 'Ресторан "Дастархан", пр. Достык 65', '12:30-14:00', 15000.00, 3),
(1, 'Посещение зоопарка', 'Алматинский зоопарк', '14:30-17:00', 2000.00, 4),
(1, 'Ужин и десерты', 'Кофейня "Sweet Family", ул. Фурманова 250', '17:30-19:00', 6000.00, 5)
ON CONFLICT DO NOTHING;

-- Шаги для корпоративного тимбилдинга
INSERT INTO steps (scenario_id, text, place, time, price, step_order) VALUES
(2, 'Утренний кофе и знакомство', 'Конференц-зал "Business Hub", ул. Розыбакиева 247', '09:00-09:30', 25000.00, 1),
(2, 'Тимбилдинг активности', 'Веревочный парк "Экстрим", ул. Навои 68', '10:00-13:00', 180000.00, 2),
(2, 'Корпоративный обед', 'Ресторан "Conference Hall", отель Rixos', '13:30-15:00', 280000.00, 3),
(2, 'Деловые игры и квесты', 'Квест-центр "Mind Games", ТРЦ Mega Park', '15:30-17:30', 120000.00, 4),
(2, 'Завершающий ужин', 'Ресторан "Panorama", БЦ Esentai Tower', '18:00-21:00', 350000.00, 5)
ON CONFLICT DO NOTHING;

-- Обновляем последовательности для корректной работы SERIAL
SELECT setval('scenarios_id_seq', (SELECT MAX(id) FROM scenarios));
SELECT setval('steps_id_seq', (SELECT MAX(id) FROM steps));