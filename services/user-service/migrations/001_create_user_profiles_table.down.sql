-- Удаляем триггер и функцию
DROP TRIGGER IF EXISTS update_user_profiles_updated_at ON user_profiles;
DROP FUNCTION IF EXISTS update_user_profiles_updated_at();

-- Удаляем индексы
DROP INDEX IF EXISTS idx_user_profiles_user_id;
DROP INDEX IF EXISTS idx_user_profiles_city;
DROP INDEX IF EXISTS idx_user_profiles_language;

-- Удаляем таблицу
DROP TABLE IF EXISTS user_profiles;