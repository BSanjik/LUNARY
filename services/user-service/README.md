# User Service - LUNARY (В разработке)

## 📝 Описание
Сервис управления профилями пользователей. Хранит персональную информацию, предпочтения и историю активности.

## ⚠️ Статус
**В разработке** - сервис еще не реализован, только заглушка

## 🚀 Планируемые эндпоинты

### GET /profile/{user_id}
Получение профиля пользователя

### PUT /profile/{user_id}
Обновление профиля пользователя

### GET /profile/{user_id}/preferences
Получение предпочтений пользователя

### PUT /profile/{user_id}/preferences
Обновление предпочтений пользователя

## 📋 Планируемая Postman Collection

### 1. Получение профиля
```json
{
  "method": "GET",
  "url": "http://localhost:8083/profile/123",
  "headers": {
    "Authorization": "Bearer {{jwt_token}}"
  }
}
```

### 2. Обновление профиля
```json
{
  "method": "PUT", 
  "url": "http://localhost:8083/profile/123",
  "headers": {
    "Authorization": "Bearer {{jwt_token}}",
    "Content-Type": "application/json"
  },
  "body": {
    "name": "Иван Петров",
    "age": 30,
    "city": "Алматы",
    "interests": ["семья", "спорт", "путешествия"]
  }
}
```

## 🗄️ Планируемая схема БД

```sql
CREATE TABLE user_profiles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER UNIQUE NOT NULL,
    name VARCHAR(255),
    age INTEGER,
    city VARCHAR(100),
    bio TEXT,
    avatar_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_preferences (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES user_profiles(user_id),
    category VARCHAR(100),
    budget_range VARCHAR(50),
    preferred_duration VARCHAR(50),
    group_size_preference VARCHAR(50),
    interests TEXT[]
);
```

## 🔧 Планируемая конфигурация

```env
DB_URL=postgres://postgres:password@postgres:5432/lunary?sslmode=disable
JWT_SECRET=secret123
PORT=8083
``` 