# 🧪 Руководство по тестированию LUNARY API

## 📋 Обзор микросервисов

| Сервис | Порт | Статус | Описание |
|--------|------|---------|----------|
| API Gateway | 8080 | ✅ Работает | Точка входа, маршрутизация |
| Auth Service | 8081 | ✅ Работает | Аутентификация и авторизация |
| Scenario Service | 8082 | ✅ Работает | Поиск сценариев событий |
| User Service | 8083 | ⚠️ В разработке | Управление профилями |
| Ad Service | 8084 | ⚠️ В разработке | Реклама и монетизация |
| Media Service | 8085 | ⚠️ В разработке | Загрузка и хранение медиа |
| Notification Service | 8086 | ⚠️ В разработке | Push-уведомления |

## 🚀 Быстрый старт

### 1. Запуск всех сервисов
```bash
docker-compose up --build
```

### 2. Проверка статуса
```bash
docker-compose ps
```

### 3. Базовый URL для API Gateway
```
http://localhost:8080
```

## 📋 Postman Collection для LUNARY

### 🔐 Переменные окружения
Создайте переменные в Postman:
- `base_url`: `http://localhost:8080`
- `auth_token`: (будет установлен автоматически после авторизации)

### 📚 Полный сценарий тестирования

## Шаг 1: Регистрация пользователя
```json
{
  "name": "1. Регистрация",
  "request": {
    "method": "POST",
    "url": "{{base_url}}/auth/registration",
    "header": [
      {
        "key": "Content-Type",
        "value": "application/json"
      }
    ],
    "body": {
      "mode": "raw",
      "raw": "{\n  \"phone\": \"+77001234567\",\n  \"password\": \"password123\",\n  \"email\": \"test@lunary.com\"\n}"
    }
  },
  "tests": "pm.test(\"Registration successful\", function () {\n    pm.response.to.have.status(201);\n    const response = pm.response.json();\n    pm.expect(response).to.have.property('token');\n    pm.environment.set('auth_token', response.token);\n});"
}
```

## Шаг 2: Авторизация пользователя  
```json
{
  "name": "2. Авторизация",
  "request": {
    "method": "POST",
    "url": "{{base_url}}/auth/login",
    "header": [
      {
        "key": "Content-Type", 
        "value": "application/json"
      }
    ],
    "body": {
      "mode": "raw",
      "raw": "{\n  \"phone\": \"+77001234567\",\n  \"password\": \"password123\"\n}"
    }
  },
  "tests": "pm.test(\"Login successful\", function () {\n    pm.response.to.have.status(200);\n    const response = pm.response.json();\n    pm.expect(response).to.have.property('token');\n    pm.environment.set('auth_token', response.token);\n});"
}
```

## Шаг 3: Поиск семейного сценария
```json
{
  "name": "3. Поиск семейного сценария",
  "request": {
    "method": "POST",
    "url": "{{base_url}}/scenario/search",
    "header": [
      {
        "key": "Content-Type",
        "value": "application/json"
      }
    ],
    "body": {
      "mode": "raw", 
      "raw": "{\n  \"query\": \"семейный ужин\"\n}"
    }
  },
  "tests": "pm.test(\"Scenario found\", function () {\n    pm.response.to.have.status(200);\n    const response = pm.response.json();\n    pm.expect(response).to.have.property('title');\n    pm.expect(response).to.have.property('checklist');\n});"
}
```

## Шаг 4: Поиск корпоративного мероприятия
```json
{
  "name": "4. Поиск корпоратива",
  "request": {
    "method": "POST",
    "url": "{{base_url}}/scenario/search",
    "header": [
      {
        "key": "Content-Type",
        "value": "application/json"
      }
    ],
    "body": {
      "mode": "raw",
      "raw": "{\n  \"query\": \"корпоратив\"\n}"
    }
  },
  "tests": "pm.test(\"Corporate scenario found\", function () {\n    pm.response.to.have.status(200);\n    const response = pm.response.json();\n    pm.expect(response.category).to.eql('corporate');\n});"
}
```

## 🧪 Тесты валидации

### Тест 1: Неверный формат телефона
```json
{
  "name": "Валидация: Неверный телефон",
  "request": {
    "method": "POST",
    "url": "{{base_url}}/auth/registration",
    "header": [{"key": "Content-Type", "value": "application/json"}],
    "body": {
      "mode": "raw",
      "raw": "{\n  \"phone\": \"invalid_phone\",\n  \"password\": \"password123\"\n}"
    }
  },
  "tests": "pm.test(\"Phone validation works\", function () {\n    pm.response.to.have.status(400);\n});"
}
```

### Тест 2: Короткий пароль
```json
{
  "name": "Валидация: Короткий пароль",
  "request": {
    "method": "POST",
    "url": "{{base_url}}/auth/registration",
    "header": [{"key": "Content-Type", "value": "application/json"}],
    "body": {
      "mode": "raw",
      "raw": "{\n  \"phone\": \"+77001234567\",\n  \"password\": \"123\"\n}"
    }
  },
  "tests": "pm.test(\"Password validation works\", function () {\n    pm.response.to.have.status(400);\n});"
}
```

### Тест 3: Пустой поисковый запрос
```json
{
  "name": "Валидация: Пустой запрос",
  "request": {
    "method": "POST",
    "url": "{{base_url}}/scenario/search",
    "header": [{"key": "Content-Type", "value": "application/json"}],
    "body": {
      "mode": "raw",
      "raw": "{\n  \"query\": \"\"\n}"
    }
  },
  "tests": "pm.test(\"Query validation works\", function () {\n    pm.response.to.have.status(400);\n});"
}
```

## 🔍 Тесты прямого доступа к сервисам

### Auth Service (прямой доступ)
```bash
# Регистрация
curl -X POST http://localhost:8081/registration \
  -H "Content-Type: application/json" \
  -d '{"phone": "+77001234567", "password": "password123"}'

# Авторизация  
curl -X POST http://localhost:8081/login \
  -H "Content-Type: application/json" \
  -d '{"phone": "+77001234567", "password": "password123"}'
```

### Scenario Service (прямой доступ)
```bash
# Поиск сценария
curl -X POST http://localhost:8082/search \
  -H "Content-Type: application/json" \
  -d '{"query": "семейный ужин"}'
```

## 📊 Ожидаемые результаты

### ✅ Успешные ответы
- **Регистрация**: 201 Created + JWT токен
- **Авторизация**: 200 OK + JWT токен  
- **Поиск сценария**: 200 OK + объект сценария

### ❌ Ошибки валидации
- **400 Bad Request**: Неверные данные
- **401 Unauthorized**: Неверные учетные данные
- **404 Not Found**: Сценарий не найден
- **409 Conflict**: Пользователь уже существует

## 🚀 Newman (CLI тестирование)

Для автоматических тестов используйте Newman:

```bash
# Установка Newman
npm install -g newman

# Экспорт коллекции из Postman в lunary-tests.json
# Запуск тестов
newman run lunary-tests.json -e lunary-environment.json
```

## 📈 Мониторинг производительности

Добавьте в тесты проверки времени ответа:

```javascript
pm.test("Response time is less than 500ms", function () {
    pm.expect(pm.response.responseTime).to.be.below(500);
});
```

## 🔧 Troubleshooting

### Сервисы не отвечают
```bash
# Проверить статус контейнеров
docker-compose ps

# Посмотреть логи
docker-compose logs api-gateway
docker-compose logs auth-service
docker-compose logs scenario-service
```

### Ошибки подключения к БД
```bash
# Проверить PostgreSQL
docker-compose logs postgres

# Перезапустить сервисы
docker-compose down -v
docker-compose up --build
``` 