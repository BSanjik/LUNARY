#!/bin/sh
set -e

echo "🔁 Ждём postgres на 5432..."
./wait-for-it.sh postgres:5432 --timeout=30 --strict -- echo "✅ Postgres доступен"

echo "⏳ Дополнительное ожидание стабилизации БД..."
sleep 5

echo "✅ База данных готова, запускаем scenario-service..."
./scenario-service 