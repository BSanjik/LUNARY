-- Создаем базу данных только если не существует
SELECT 'CREATE DATABASE lunary'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'lunary')\gexec