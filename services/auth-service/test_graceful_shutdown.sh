#!/bin/bash

# Скрипт для тестирования graceful shutdown
echo "🚀 Testing Graceful Shutdown..."

# Запускаем сервис в фоне
echo "📡 Starting auth service..."
go run ./cmd/main.go &
SERVER_PID=$!

# Ждем запуска сервера
echo "⏳ Waiting for server to start..."
sleep 3

# Проверяем что сервер запустился
echo "🔍 Checking if server is running..."
if curl -s http://localhost:8081/health/live > /dev/null; then
    echo "✅ Server is running"
else
    echo "❌ Server failed to start"
    kill $SERVER_PID 2>/dev/null
    exit 1
fi

# Создаем активные соединения (имитация нагрузки)
echo "📊 Creating active connections..."
for i in {1..3}; do
    curl -s http://localhost:8081/health &
done

# Отправляем SIGTERM для graceful shutdown
echo "🛑 Sending SIGTERM for graceful shutdown..."
kill -TERM $SERVER_PID

# Ждем завершения
echo "⏳ Waiting for graceful shutdown..."
wait $SERVER_PID

echo "✅ Graceful shutdown test completed!"
echo ""
echo "📝 Check the logs above for:"
echo "   - 'Received shutdown signal'"
echo "   - 'Starting graceful shutdown'"
echo "   - 'HTTP server stopped gracefully'"
echo "   - 'Cleanup function completed'"
echo "   - 'Graceful shutdown completed successfully'"
