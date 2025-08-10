#!/bin/bash

# Скрипт для тестирования health checks
echo "🚀 Testing Auth Service Health Checks..."

BASE_URL="http://localhost:8081"

echo ""
echo "📊 Testing /health endpoint..."
curl -s "$BASE_URL/health" | jq '.'

echo ""
echo "❤️ Testing /health/live endpoint..."
curl -s "$BASE_URL/health/live" | jq '.'

echo ""
echo "✅ Testing /health/ready endpoint..."
curl -s "$BASE_URL/health/ready" | jq '.'

echo ""
echo "📈 Health check testing completed!"
