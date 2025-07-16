# Scenario Service - LUNARY

## 📝 Описание
Сервис для поиска и управления сценариями событий. Предоставляет каталог готовых сценариев для различных мероприятий (семейные встречи, корпоративы, дружеские вечера).

## 🚀 Запуск
```bash
# В корне проекта
docker-compose up scenario-service

# Или локально
cd services/scenario-service
go run ./cmd
```

## 🗄️ База данных
Использует PostgreSQL с таблицами:

### Таблица `scenarios`:
```sql
CREATE TABLE scenarios (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100),
    duration VARCHAR(50),
    participants VARCHAR(50),
    budget VARCHAR(50),
    checklist TEXT[],
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Таблица `feedback`:
```sql
CREATE TABLE feedback (
    id SERIAL PRIMARY KEY,
    scenario_id INTEGER REFERENCES scenarios(id),
    user_login VARCHAR(255),
    rating INTEGER CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 🌐 Эндпоинты

### POST /search
Поиск сценариев по ключевым словам

**Параметры:**
- `query` (string) - Поисковый запрос

## 📋 Postman Collection

### 1. Поиск семейного сценария
```json
{
  "method": "POST",
  "url": "http://localhost:8082/search",
  "headers": {
    "Content-Type": "application/json"
  },
  "body": {
    "query": "семейный ужин"
  }
}
```

**Ожидаемый ответ (200):**
```json
{
  "id": 1,
  "title": "Семейный ужин дома",
  "description": "Уютный вечер с семьей за домашним столом. Время для общения, вкусной еды и теплой атмосферы.",
  "category": "family",
  "duration": "2-3 часа",
  "participants": "3-6 человек",
  "budget": "средний",
  "checklist": [
    "Составить меню на основе предпочтений семьи",
    "Купить свежие продукты",
    "Приготовить основные блюда и закуски",
    "Сервировать стол с красивой посудой",
    "Подготовить фоновую музыку",
    "Придумать темы для разговоров"
  ]
}
```

### 2. Поиск корпоративного мероприятия
```json
{
  "method": "POST",
  "url": "http://localhost:8082/search",
  "headers": {
    "Content-Type": "application/json"
  },
  "body": {
    "query": "корпоратив"
  }
}
```

**Ожидаемый ответ (200):**
```json
{
  "id": 2,
  "title": "Корпоративная вечеринка",
  "description": "Незабываемый корпоративный праздник для сплочения команды и празднования достижений.",
  "category": "corporate",
  "duration": "4-5 часов",
  "participants": "15-50 человек",
  "budget": "высокий",
  "checklist": [
    "Забронировать подходящее помещение",
    "Составить развлекательную программу",
    "Организовать кейтеринг или банкет",
    "Подготовить корпоративные подарки",
    "Продумать музыкальное сопровождение",
    "Назначить ответственных за различные части мероприятия"
  ]
}
```

### 3. Поиск несуществующего сценария
```json
{
  "method": "POST",
  "url": "http://localhost:8082/search",
  "headers": {
    "Content-Type": "application/json"
  },
  "body": {
    "query": "космическое путешествие"
  }
}
```

**Ожидаемый ответ (404):**
```json
{
  "error": "Сценарий не найден. Попробуйте изменить запрос"
}
```

### 4. Тест с пустым запросом
```json
{
  "method": "POST",
  "url": "http://localhost:8082/search",
  "headers": {
    "Content-Type": "application/json"
  },
  "body": {
    "query": ""
  }
}
```

**Ожидаемый ответ (400):**
```json
{
  "error": "Поле 'query' обязательно для заполнения"
}
```

### 5. Тест с неверным JSON
```json
{
  "method": "POST",
  "url": "http://localhost:8082/search",
  "headers": {
    "Content-Type": "application/json"
  },
  "body": "invalid json"
}
```

**Ожидаемый ответ (400):**
```json
{
  "error": "Некорректный формат запроса"
}
```

## 📊 Предустановленные сценарии

### 1. Семейные мероприятия
- Семейный ужин дома
- Детский день рождения
- Семейный пикник

### 2. Корпоративные события  
- Корпоративная вечеринка
- Тимбилдинг на природе
- Презентация продукта

### 3. Дружеские встречи
- Вечер настольных игр
- Домашняя вечеринка
- Барбекю на даче

## 🔧 Конфигурация

Переменные окружения (`.env`):
```env
DB_URL=postgres://postgres:password@postgres:5432/lunary?sslmode=disable
PORT=8082
```

## 🏗️ Архитектура
```
cmd/main.go          # Точка входа
├── config/          # Конфигурация
├── storage/         # Работа с БД
├── service/         # Бизнес-логика поиска
├── handler/         # HTTP обработчики
└── model/           # Модели данных
```

## 🔍 Алгоритм поиска
1. Поиск точного совпадения по заголовку
2. Поиск по ключевым словам в описании
3. Поиск по категории
4. Возврат первого найденного результата

## 🚀 Будущие улучшения
- [ ] Возврат множественных результатов поиска
- [ ] Фильтрация по категориям, бюджету, количеству участников
- [ ] Система рейтингов и отзывов
- [ ] AI-генерация персонализированных сценариев 