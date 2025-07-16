# API Спецификация LUNARY

## Обзор

LUNARY предоставляет RESTful API для всех операций платформы планирования мероприятий. API построен на принципах микросервисной архитектуры, где каждый сервис отвечает за свою область функциональности.

## Базовая информация

- **Base URL**: `http://localhost:8080` (через API Gateway)
- **Protocol**: HTTP/1.1 (HTTPS в production)
- **Content-Type**: `application/json`
- **Authentication**: JWT Bearer Token
- **API Version**: v1

## Аутентификация

Все защищенные эндпоинты требуют JWT токен в заголовке:

```http
Authorization: Bearer <jwt_token>
```

## API Gateway (Port 8080)

API Gateway является единой точкой входа для всех запросов и маршрутизирует их к соответствующим микросервисам.

### Health Check

#### GET /health
Проверка работоспособности API Gateway.

**Ответ:**
```json
{
  "status": "ok",
  "timestamp": "2024-01-15T10:30:00Z",
  "services": {
    "auth-service": "healthy",
    "scenario-service": "healthy"
  }
}
```

## Auth Service API

### Регистрация

#### POST /auth/register
Регистрация нового пользователя.

**Параметры запроса:**
```json
{
  "name": "string (required, min: 2, max: 100)",
  "phone": "string (required, format: digits only, min: 10, max: 15)",
  "password": "string (required, min: 6, max: 100)"
}
```

**Пример запроса:**
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Иван Иванов",
    "phone": "79123456789",
    "password": "securepassword123"
  }'
```

**Успешный ответ (201 Created):**
```json
{
  "message": "User created successfully",
  "user": {
    "id": 1,
    "name": "Иван Иванов",
    "phone": "79123456789",
    "created_at": "2024-01-15T10:30:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Ошибки:**
- `400 Bad Request` - Некорректные данные
- `409 Conflict` - Пользователь уже существует

### Авторизация

#### POST /auth/login
Авторизация пользователя.

**Параметры запроса:**
```json
{
  "phone": "string (required)",
  "password": "string (required)"
}
```

**Пример запроса:**
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "79123456789",
    "password": "securepassword123"
  }'
```

**Успешный ответ (200 OK):**
```json
{
  "message": "Login successful",
  "user": {
    "id": 1,
    "name": "Иван Иванов",
    "phone": "79123456789"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Ошибки:**
- `400 Bad Request` - Некорректные данные
- `401 Unauthorized` - Неверные учетные данные

### Проверка токена

#### GET /auth/validate
Проверка валидности JWT токена.

**Заголовки:**
```http
Authorization: Bearer <jwt_token>
```

**Успешный ответ (200 OK):**
```json
{
  "valid": true,
  "user": {
    "id": 1,
    "name": "Иван Иванов",
    "phone": "79123456789"
  }
}
```

**Ошибки:**
- `401 Unauthorized` - Недействительный или отсутствующий токен

## Scenario Service API

### Поиск сценариев

#### GET /scenarios/search
Поиск сценариев мероприятий по ключевым словам.

**Заголовки:**
```http
Authorization: Bearer <jwt_token>
```

**Query параметры:**
- `query` (string, required) - Поисковый запрос
- `limit` (integer, optional, default: 10, max: 50) - Количество результатов
- `offset` (integer, optional, default: 0) - Смещение для пагинации

**Пример запроса:**
```bash
curl -X GET "http://localhost:8080/scenarios/search?query=день рождения&limit=5" \
  -H "Authorization: Bearer <jwt_token>"
```

**Успешный ответ (200 OK):**
```json
{
  "scenarios": [
    {
      "id": 1,
      "title": "День рождения в стиле пиратов",
      "description": "Захватывающий пиратский праздник для детей от 5 до 12 лет",
      "category": "Детский праздник",
      "duration": "3 часа",
      "participants": "5-15 человек",
      "budget": "10000-25000 руб",
      "tags": ["дети", "пираты", "активности", "костюмы"],
      "rating": 4.8,
      "created_at": "2024-01-10T15:30:00Z"
    },
    {
      "id": 2,
      "title": "Взрослый день рождения в ресторане",
      "description": "Элегантное празднование дня рождения в ресторане",
      "category": "Взрослый праздник",
      "duration": "4 часа",
      "participants": "10-50 человек",
      "budget": "50000-150000 руб",
      "tags": ["взрослые", "ресторан", "торжественно"],
      "rating": 4.5,
      "created_at": "2024-01-12T10:15:00Z"
    }
  ],
  "total": 15,
  "limit": 5,
  "offset": 0
}
```

**Ошибки:**
- `400 Bad Request` - Некорректные параметры запроса
- `401 Unauthorized` - Отсутствует или недействительный токен

### Получение сценария по ID

#### GET /scenarios/{id}
Получение подробной информации о конкретном сценарии.

**Заголовки:**
```http
Authorization: Bearer <jwt_token>
```

**Path параметры:**
- `id` (integer, required) - ID сценария

**Пример запроса:**
```bash
curl -X GET http://localhost:8080/scenarios/1 \
  -H "Authorization: Bearer <jwt_token>"
```

**Успешный ответ (200 OK):**
```json
{
  "id": 1,
  "title": "День рождения в стиле пиратов",
  "description": "Захватывающий пиратский праздник для детей от 5 до 12 лет",
  "detailed_description": "Полное описание сценария с пошаговым планом...",
  "category": "Детский праздник",
  "duration": "3 часа",
  "participants": "5-15 человек",
  "budget": "10000-25000 руб",
  "tags": ["дети", "пираты", "активности", "костюмы"],
  "rating": 4.8,
  "materials": [
    "Пиратские костюмы",
    "Карта сокровищ",
    "Сундук с сокровищами",
    "Пиратский флаг"
  ],
  "steps": [
    {
      "step": 1,
      "title": "Подготовка пространства",
      "description": "Украшение помещения в пиратском стиле",
      "duration": "30 минут"
    },
    {
      "step": 2,
      "title": "Приветствие и раздача костюмов",
      "description": "Встреча гостей и превращение в пиратов",
      "duration": "15 минут"
    }
  ],
  "created_at": "2024-01-10T15:30:00Z"
}
```

**Ошибки:**
- `404 Not Found` - Сценарий не найден
- `401 Unauthorized` - Отсутствует или недействительный токен

### Отправка обратной связи

#### POST /scenarios/{id}/feedback
Отправка отзыва о сценарии.

**Заголовки:**
```http
Authorization: Bearer <jwt_token>
```

**Path параметры:**
- `id` (integer, required) - ID сценария

**Параметры запроса:**
```json
{
  "rating": "integer (required, min: 1, max: 5)",
  "comment": "string (optional, max: 1000)"
}
```

**Пример запроса:**
```bash
curl -X POST http://localhost:8080/scenarios/1/feedback \
  -H "Authorization: Bearer <jwt_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "rating": 5,
    "comment": "Отличный сценарий! Дети были в восторге!"
  }'
```

**Успешный ответ (201 Created):**
```json
{
  "message": "Feedback submitted successfully",
  "feedback": {
    "id": 1,
    "scenario_id": 1,
    "user_id": 1,
    "rating": 5,
    "comment": "Отличный сценарий! Дети были в восторге!",
    "created_at": "2024-01-15T14:20:00Z"
  }
}
```

**Ошибки:**
- `400 Bad Request` - Некорректные данные
- `404 Not Found` - Сценарий не найден
- `401 Unauthorized` - Отсутствует или недействительный токен
- `409 Conflict` - Отзыв уже оставлен

## User Service API (В разработке)

### Получение профиля

#### GET /users/profile
Получение профиля текущего пользователя.

**Заголовки:**
```http
Authorization: Bearer <jwt_token>
```

**Планируемый ответ (200 OK):**
```json
{
  "user": {
    "id": 1,
    "name": "Иван Иванов",
    "phone": "79123456789",
    "email": "ivan@example.com",
    "preferences": {
      "favorite_categories": ["Детский праздник", "Семейный отдых"],
      "budget_range": "10000-50000",
      "notifications_enabled": true
    },
    "statistics": {
      "scenarios_viewed": 25,
      "feedback_given": 5,
      "events_planned": 3
    },
    "created_at": "2024-01-10T10:00:00Z",
    "updated_at": "2024-01-15T14:20:00Z"
  }
}
```

### Обновление профиля

#### PUT /users/profile
Обновление профиля пользователя.

**Заголовки:**
```http
Authorization: Bearer <jwt_token>
```

**Параметры запроса:**
```json
{
  "name": "string (optional)",
  "email": "string (optional, format: email)",
  "preferences": {
    "favorite_categories": ["string"],
    "budget_range": "string",
    "notifications_enabled": "boolean"
  }
}
```

## Ad Service API (Планируется)

### Создание объявления

#### POST /ads
Создание нового рекламного объявления.

**Заголовки:**
```http
Authorization: Bearer <jwt_token>
```

**Планируемые параметры:**
```json
{
  "title": "string (required)",
  "description": "string (required)",
  "category": "string (required)",
  "price": "number (optional)",
  "location": "string (optional)",
  "contact_info": {
    "phone": "string",
    "email": "string"
  },
  "images": ["string"] // URLs изображений
}
```

### Поиск объявлений

#### GET /ads/search
Поиск рекламных объявлений.

**Query параметры:**
- `query` - поисковый запрос
- `category` - категория
- `location` - местоположение
- `min_price`, `max_price` - диапазон цен

## Media Service API (Планируется)

### Загрузка файла

#### POST /media/upload
Загрузка медиафайла.

**Заголовки:**
```http
Authorization: Bearer <jwt_token>
Content-Type: multipart/form-data
```

**Планируемый ответ:**
```json
{
  "file": {
    "id": "uuid",
    "filename": "original_name.jpg",
    "url": "https://media.lunary.app/files/uuid.jpg",
    "size": 1024000,
    "content_type": "image/jpeg",
    "uploaded_at": "2024-01-15T15:30:00Z"
  }
}
```

## Notification Service API (Планируется)

### Отправка уведомления

#### POST /notifications/send
Отправка push-уведомления пользователю.

**Заголовки:**
```http
Authorization: Bearer <jwt_token>
```

**Планируемые параметры:**
```json
{
  "user_id": "integer",
  "title": "string",
  "message": "string",
  "type": "push | email | sms",
  "data": {} // Дополнительные данные
}
```

## Коды ошибок

### Стандартные HTTP коды

- `200 OK` - Успешный запрос
- `201 Created` - Ресурс создан
- `400 Bad Request` - Некорректные данные
- `401 Unauthorized` - Не авторизован
- `403 Forbidden` - Доступ запрещен
- `404 Not Found` - Ресурс не найден
- `409 Conflict` - Конфликт данных
- `422 Unprocessable Entity` - Ошибка валидации
- `429 Too Many Requests` - Превышен лимит запросов
- `500 Internal Server Error` - Внутренняя ошибка сервера

### Формат ошибок

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": [
      {
        "field": "phone",
        "message": "Phone number must contain only digits"
      }
    ]
  }
}
```

### Коды пользовательских ошибок

- `USER_ALREADY_EXISTS` - Пользователь уже существует
- `INVALID_CREDENTIALS` - Неверные учетные данные
- `TOKEN_EXPIRED` - Токен истек
- `VALIDATION_ERROR` - Ошибка валидации
- `SCENARIO_NOT_FOUND` - Сценарий не найден
- `RATE_LIMIT_EXCEEDED` - Превышен лимит запросов

## Rate Limiting

API Gateway ограничивает количество запросов:

- **Лимит**: 100 запросов в минуту на IP
- **Заголовки ответа**:
  - `X-RateLimit-Limit: 100`
  - `X-RateLimit-Remaining: 95`
  - `X-RateLimit-Reset: 1642248000`

При превышении лимита возвращается:
```json
{
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Too many requests"
  }
}
```

## Пагинация

Для эндпоинтов, возвращающих списки, используется offset-based пагинация:

**Параметры:**
- `limit` - количество элементов (max: 50, default: 10)
- `offset` - смещение (default: 0)

**Ответ:**
```json
{
  "data": [...],
  "total": 150,
  "limit": 10,
  "offset": 20,
  "has_more": true
}
```

## Версионирование

В будущем планируется добавить версионирование API:

```http
Accept: application/vnd.lunary.v1+json
```

или через URL:
```
/api/v1/scenarios/search
```

## CORS

API Gateway настроен для CORS:

```
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: Content-Type, Authorization
```

## Swagger документация

Интерактивная документация доступна по адресам:

- **Auth Service**: http://localhost:8081/swagger/index.html
- **Scenario Service**: http://localhost:8082/swagger/index.html (планируется)

## Примеры использования

### Полный flow регистрации и поиска

```bash
# 1. Регистрация
RESPONSE=$(curl -s -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","phone":"1234567890","password":"password123"}')

# 2. Извлечение токена
TOKEN=$(echo $RESPONSE | jq -r '.token')

# 3. Поиск сценариев
curl -X GET "http://localhost:8080/scenarios/search?query=день рождения" \
  -H "Authorization: Bearer $TOKEN"

# 4. Получение конкретного сценария
curl -X GET http://localhost:8080/scenarios/1 \
  -H "Authorization: Bearer $TOKEN"

# 5. Отправка отзыва
curl -X POST http://localhost:8080/scenarios/1/feedback \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"rating":5,"comment":"Отличный сценарий!"}'
```

## Тестирование

Полная коллекция Postman доступна в [TESTING.md](../TESTING.md).

Для автоматического тестирования используйте Newman:
```bash
newman run LUNARY_API_Tests.postman_collection.json
```