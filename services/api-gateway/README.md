# API Gateway - LUNARY

## 📝 Описание
API Gateway является точкой входа для всех HTTP запросов в микросервисную архитектуру LUNARY. Он маршрутизирует запросы к соответствующим сервисам и обеспечивает аутентификацию.

## 🚀 Запуск
```bash
# В корне проекта
docker-compose up api-gateway

# Или локально
cd services/api-gateway
go run ./cmd
```

## 🌐 Эндпоинты

### Аутентификация (проксирует в auth-service)
- `POST /auth/registration` - Регистрация пользователя
- `POST /auth/login` - Авторизация пользователя  
- `GET /auth/docs/` - Swagger документация

### Сценарии (проксирует в scenario-service)
- `POST /scenario/search` - Поиск сценариев

## 📋 Postman Collection

### 1. Регистрация пользователя
```json
{
  "method": "POST",
  "url": "http://localhost:8080/auth/registration",
  "headers": {
    "Content-Type": "application/json"
  },
  "body": {
    "phone": "+77001234567",
    "password": "password123", 
    "email": "test@example.com"
  }
}
```

**Ожидаемый ответ:**
```json
{
  "message": "Пользователь успешно зарегистрирован",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 2. Авторизация пользователя
```json
{
  "method": "POST",
  "url": "http://localhost:8080/auth/login",
  "headers": {
    "Content-Type": "application/json"
  },
  "body": {
    "phone": "+77001234567",
    "password": "password123"
  }
}
```

**Ожидаемый ответ:**
```json
{
  "message": "Авторизация успешна",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 3. Поиск сценариев
```json
{
  "method": "POST", 
  "url": "http://localhost:8080/scenario/search",
  "headers": {
    "Content-Type": "application/json"
  },
  "body": {
    "query": "семейный ужин"
  }
}
```

**Ожидаемый ответ:**
```json
{
  "id": 1,
  "title": "Семейный ужин дома",
  "description": "Уютный вечер с семьей за домашним столом...",
  "category": "family",
  "duration": "2-3 часа",
  "participants": "3-6 человек",
  "budget": "средний",
  "checklist": [
    "Составить меню",
    "Купить продукты", 
    "Приготовить блюда"
  ]
}
```

## 🔧 Конфигурация

Переменные окружения (`.env`):
```env
JWT_SECRET=secret123
LISTEN_ADDR=:8080
SERVICE_AUTH=http://auth-service:8081
SERVICE_SCENARIO=http://scenario-service:8082
```

## 📊 Мониторинг
- Gateway работает на порту `8080`
- Логи маршрутизации запросов в stdout
- Rate limiting: 100 запросов в минуту на IP 