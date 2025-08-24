-- Создаем таблицу профилей с НАСТОЯЩИМ Foreign Key
CREATE TABLE IF NOT EXISTS user_profiles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL UNIQUE 
        REFERENCES users(id) ON DELETE CASCADE,  -- ✅ Настоящий FK!
    name VARCHAR(50),
    lastname VARCHAR(50),
    sex VARCHAR(10) CHECK (sex IN ('MALE', 'FEMALE')),
    profession VARCHAR(100),
    city VARCHAR(100) NOT NULL DEFAULT 'Almaty',
    timezone VARCHAR(50) NOT NULL DEFAULT 'Asia/Almaty',
    language VARCHAR(5) NOT NULL DEFAULT 'RU' 
        CHECK (language IN ('RU', 'KZ', 'EN')),
    notification_tone VARCHAR(20) NOT NULL DEFAULT 'NEUTRAL' 
        CHECK (notification_tone IN ('SOFT', 'NEUTRAL', 'STRICT')),
    last_special_scenario_date TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Уникальный индекс для связи 1:1
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_profiles_user_id ON user_profiles(user_id);

-- Остальные индексы...
CREATE INDEX IF NOT EXISTS idx_user_profiles_city ON user_profiles(city);
CREATE INDEX IF NOT EXISTS idx_user_profiles_language ON user_profiles(language);

-- Триггер для updated_at
CREATE OR REPLACE FUNCTION update_user_profiles_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_user_profiles_updated_at 
    BEFORE UPDATE ON user_profiles 
    FOR EACH ROW EXECUTE FUNCTION update_user_profiles_updated_at();

    -- Комментарии для документации
COMMENT ON TABLE user_profiles IS 'Профили пользователей с FK на users.id';
COMMENT ON COLUMN user_profiles.user_id IS 'Foreign Key на users.id с каскадным удалением';