# Auth Service - LUNARY

## 📝 Описание
Сервис аутентификации и авторизации LUNARY. Обрабатывает регистрацию пользователей, авторизацию и выдачу JWT токенов.

## 🚀 Запуск
```bash
# В корне проекта
docker-compose up auth-service

# Или локально
cd services/auth-service
go run ./cmd
```

## 🗄️ База данных
Использует PostgreSQL с таблицей `users`:
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    phone VARCHAR(20) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,  
    email VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 🌐 Эндпоинты

### POST /registration
Регистрация нового пользователя

**Параметры:**
- `phone` (string) - Номер телефона в формате +7XXXXXXXXXX
- `password` (string) - Пароль минимум 6 символов
- `email` (string, optional) - Email адрес

### POST /login  
Авторизация пользователя

**Параметры:**
- `phone` (string) - Номер телефона
- `password` (string) - Пароль

### GET /docs/
Swagger документация API

## 📋 Postman Collection

### 1. Регистрация пользователя
```json
{
  "method": "POST",
  "url": "http://localhost:8081/registration",
  "headers": {
    "Content-Type": "application/json"
  },
  "body": {
    "phone": "+77001234567",
    "password": "password123",
    "email": "user@example.com"
  }
}
```

**Успешный ответ (201):**
```json
{
  "message": "Пользователь успешно зарегистрирован",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwaG9uZSI6Iis3NzAwMTIzNDU2NyIsImV4cCI6MTYzOTU2NzgwMH0.signature"
}
```

**Ошибки:**
- `400` - Неверный формат телефона или короткий пароль
- `409` - Пользователь с таким телефоном уже существует

### 2. Авторизация пользователя
```json
{
  "method": "POST", 
  "url": "http://localhost:8081/login",
  "headers": {
    "Content-Type": "application/json"
  },
  "body": {
    "phone": "+77001234567",
    "password": "password123"
  }
}
```

**Успешный ответ (200):**
```json
{
  "message": "Авторизация успешна",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwaG9uZSI6Iis3NzAwMTIzNDU2NyIsImV4cCI6MTYzOTU2NzgwMH0.signature"
}
```

**Ошибки:**
- `400` - Неверный формат телефона
- `401` - Неверный пароль или пользователь не найден

### 3. Тест валидации телефона
```json
{
  "method": "POST",
  "url": "http://localhost:8081/registration", 
  "headers": {
    "Content-Type": "application/json"
  },
  "body": {
    "phone": "invalid_phone",
    "password": "password123"
  }
}
```

**Ожидаемый ответ (400):**
```json
{
  "error": "Неверный формат телефона"
}
```

### 4. Тест короткого пароля
```json
{
  "method": "POST",
  "url": "http://localhost:8081/registration",
  "headers": {
    "Content-Type": "application/json"
  }, 
  "body": {
    "phone": "+77001234567",
    "password": "123"
  }
}
```

**Ожидаемый ответ (400):**
```json
{
  "error": "Пароль должен содержать минимум 6 символов"
}
```

## 🔧 Конфигурация

Переменные окружения (`.env`):
```env
DB_URL=postgres://postgres:password@postgres:5432/lunary?sslmode=disable
JWT_SECRET=secret123
PORT=8081
```

## 🔐 Безопасность
- Пароли хэшируются с использованием bcrypt
- JWT токены с истечением срока действия (24 часа)
- Валидация формата телефона regex
- Защита от SQL-инъекций через prepared statements

## 🏗️ Архитектура
```
cmd/main.go          # Точка входа
├── config/          # Конфигурация
├── storage/         # Работа с БД
├── service/         # Бизнес-логика  
├── handler/         # HTTP обработчики
└── token/           # JWT токены
``` 