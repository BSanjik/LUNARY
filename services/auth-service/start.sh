#!/bin/sh
set -e

export PATH=$PATH:/usr/local/bin

echo "🔁 Ждём postgres на 5432..."
./wait-for-it.sh postgres:5432 --timeout=30 --strict -- echo "✅ Postgres доступен"

echo "🚀 Выполняем миграции..."
migrate -path=./migrations -database "postgres://postgres:mypassword@postgres:5432/lunary_auth?sslmode=disable" up

echo "✅ Миграции завершены, запускаем сервис..."
./auth-service
