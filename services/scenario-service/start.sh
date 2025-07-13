#!/bin/bash
set -e

export PATH=$PATH:/usr/local/bin

echo "🔁 Ждём postgres на 5432..."
./wait-for-it.sh postgres:5432 --timeout=30 --strict -- echo "✅ Postgres доступен"

echo "🚀 Выполняем миграции..."
migrate -path=./migrations -database "${DB_URL}" up

echo "✅ Миграции завершены, запускаем сервис..."
exec ./scenario-service