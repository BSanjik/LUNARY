# Техническое задание: Lunary AI Assistant (Обновленная версия)

## 1. Общее описание проекта

### 1.1 Назначение
**Lunary AI Assistant** — интеллектуальный чат-ассистент для персонального планирования, интегрированный в Flutter приложение. Использует AI (Google Gemini) для создания индивидуальных планов дня и специальных сценариев на основе пользовательских анкет через чат-интерфейс.

### 1.2 Ключевые изменения от исходного ТЗ
- ❌ **Убрали**: Telegram Bot API, webhooks, внешние боты
- ✅ **Заменили на**: Встроенный AI-чат в Flutter приложении
- ✅ **Упростили**: Одна база данных (PostgreSQL с JSONB)
- ✅ **Интегрировали**: С существующей микросервисной архитектурой

### 1.3 Цели проекта
- Автоматизация процесса планирования дня через чат
- Персонализированные рекомендации на основе AI
- Умные напоминания и мотивационные сообщения
- Создание специальных сценариев для событий
- Нативный опыт внутри приложения

### 1.4 Целевая аудитория
- Пользователи Flutter приложения LUNARY
- Люди, желающие оптимизировать свой день
- Пользователи, нуждающиеся в мотивации и структуре
- Активные люди, планирующие специальные события

## 2. Архитектура системы

### 2.1 Обновленная микросервисная архитектура

```
Flutter App (Chat UI) → API Gateway → AI Assistant Service ← AI Service
                                          ↓                    ↓
                      Auth Service ← PostgreSQL → Plan Service
                         ↓                           ↓
                   User Service ← → Notification Service
```

### 2.2 Сервисы

#### 2.2.1 Существующие сервисы (расширяем)
| Сервис | Порт | Статус | Новый функционал |
|--------|------|--------|------------------|
| **API Gateway** | 8080 | ✅ Готов | Роуты для чата |
| **Auth Service** | 8081 | ✅ Готов | Расширенные профили пользователей |
| **AI Assistant Service** | 8082 | 🔄 Трансформация | Чат + AI логика (ex-Scenario Service) |

#### 2.2.2 Новые сервисы
| Сервис | Порт | Статус | Описание |
|--------|------|--------|----------|
| **AI Service** | 8088 | 🆕 Новый | Google Gemini интеграция |
| **Plan Service** | 8089 | 🆕 Новый | Хранение планов и анкет |

#### 2.2.3 Планируемые сервисы (дорабатываем)
| Сервис | Порт | Статус | Функционал |
|--------|------|--------|------------|
| **User Service** | 8083 | 🔄 В разработке | Профили с полями из ТЗ |
| **Notification Service** | 8086 | 🔄 В разработке | Push-уведомления для планов |

### 2.3 База данных: PostgreSQL ONLY

**Почему PostgreSQL, а не MongoDB:**
- ✅ **JSONB поддержка** - все возможности NoSQL
- ✅ **ACID транзакции** - надежность данных
- ✅ **Меньше сложности** - одна БД вместо двух
- ✅ **Команда знает PostgreSQL**
- ✅ **Экономия ресурсов**

## 3. Функциональные требования

### 3.1 Система чата с AI ассистентом

#### 3.1.1 Chat API
**Основные endpoints:**
```
POST /chat/message           - отправка сообщения
GET  /chat/history/{user_id} - история чата
POST /chat/clear/{user_id}   - очистка истории
```

**Модель чат-сообщения:**
```json
{
  "id": "uuid",
  "user_id": 123,
  "message": "Помоги составить план на завтра",
  "response": "Конечно! Начнем с вечерней анкеты...",
  "message_type": "text|survey|scenario|plan",
  "metadata": {
    "buttons": [],
    "state": "evening_survey_start"
  },
  "created_at": "2025-01-15T10:30:00Z"
}
```

#### 3.1.2 State Machine для чата
**Состояния пользователя:**
```
start → main_menu → {
  evening_survey_start → waiting_wake_time → waiting_sleep_time → ...
  special_scenario_start → waiting_scenario_type → ...
  viewing_plans → plan_details → ...
}
```

### 3.2 Система регистрации и профилей (расширяем Auth Service)

#### 3.2.1 Расширенная модель пользователя
```sql
-- Добавляем к существующей таблице users
ALTER TABLE users ADD COLUMN name VARCHAR(50);
ALTER TABLE users ADD COLUMN lastname VARCHAR(50);
ALTER TABLE users ADD COLUMN sex VARCHAR(10) CHECK (sex IN ('MALE', 'FEMALE'));
ALTER TABLE users ADD COLUMN profession VARCHAR(100);
ALTER TABLE users ADD COLUMN city VARCHAR(100) DEFAULT 'Almaty';
ALTER TABLE users ADD COLUMN timezone VARCHAR(50) DEFAULT 'Asia/Almaty';
ALTER TABLE users ADD COLUMN language VARCHAR(5) DEFAULT 'RU';
ALTER TABLE users ADD COLUMN notification_tone VARCHAR(20) DEFAULT 'NEUTRAL';
ALTER TABLE users ADD COLUMN last_special_scenario_date TIMESTAMP;
```

### 3.3 Система планирования дня

#### 3.3.1 Вечерняя анкета через чат
**Поля анкеты (собираются через чат):**
- Время пробуждения на завтра (HH:MM)
- Время отбоя на завтра (HH:MM)
- Главная цель дня (обязательно)
- Уровень энергии (1-5)
- Стиль дня (Спокойный/Динамичный/Фокус на работе/Семейный)
- Обязательные дела/встречи (опционально)
- Второстепенные цели (опционально)
- Ограничения (опционально)

**Пример диалога в чате:**
```
🤖: Привет! Давай составим план на завтра. Во сколько планируешь проснуться?
👤: В 8 утра
🤖: Понял: 08:00. А во сколько лечь спать?
👤: В 23:00
🤖: Отлично! Какая главная цель завтрашнего дня?
👤: Закончить проект и сходить в спортзал
...
```

#### 3.3.2 Генерация ежедневного плана
**Через AI Service → Google Gemini:**
```json
{
  "plan_text": "Динамичный день с фокусом на проект",
  "timeline": [
    {
      "start": "08:30",
      "end": "09:00", 
      "title": "Утренняя рутина",
      "text": "Завтрак и зарядка для энергии",
      "nudge_text": "Время начинать день!",
      "tags": ["morning", "energy"]
    }
  ],
  "motivation_texts": {
    "morning": "Утренняя мотивация (≤160 символов)",
    "day": "Дневная мотивация (≤160 символов)", 
    "evening": "Вечерняя мотивация (≤160 символов)"
  }
}
```

### 3.4 Система специальных сценариев

#### 3.4.1 Анкета специального сценария через чат
**Поля (собираются в диалоге):**
- Тип сценария (Date/Family/Friends/Other)
- Город/район проведения
- Дата и время начала
- Длительность (30-1440 минут)
- Количество участников и состав
- Бюджет (числовое значение)
- Способ передвижения (walking/public/car/taxi/mixed)
- Настроение/вайб (список через запятую)
- Цель сценария
- Ограничения (опционально)

#### 3.4.2 Генерация вариантов сценария
**AI генерирует варианты:**
```json
{
  "options": [
    {
      "title": "Романтический вечер в центре",
      "summary": "Уютное свидание с ужином и прогулкой",
      "timeline": [
        {
          "at_rel_min": 0,
          "title": "Встреча в кафе",
          "text": "Кофе и легкий разговор"
        }
      ],
      "checklist": ["забронировать столик", "выбрать наряд"],
      "vibe": ["romantic", "cozy"],
      "est_cost": "mid",
      "indoor_outdoor": "mix",
      "planB": "Кино при непогоде"
    }
  ]
}
```

### 3.5 Система просмотра планов через чат

#### 3.5.1 Команды чата для планов
```
👤: Покажи мои планы
🤖: У тебя есть планы на:
    📅 15 января - Ежедневный план
    ✨ 18 января - Романтический вечер
    
    Какой хочешь посмотреть?

👤: 15 января
🤖: 📅 План на 15 января:
    08:30-09:00 Утренняя рутина
    09:00-12:00 Работа над проектом
    ...
```

## 4. Технические требования

### 4.1 Технологический стек

#### 4.1.1 Основные технологии
- **Язык:** Go 1.21+
- **HTTP Framework:** Gin (как в существующих сервисах)
- **AI:** Google Generative AI SDK for Go
- **База данных:** PostgreSQL 15+ с JSONB
- **ORM:** GORM (существующий)
- **Миграции:** golang-migrate/migrate (существующий)
- **Контейнеризация:** Docker + Docker Compose (готово)

#### 4.1.2 Новые зависимости
```go
// go.mod additions
github.com/google/generative-ai-go v0.5.0
github.com/google/uuid v1.3.0  // для chat message IDs
```

### 4.2 Модели данных

#### 4.2.1 Chat Messages (PostgreSQL)
```sql
CREATE TABLE chat_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    response TEXT NOT NULL,
    message_type VARCHAR(50) DEFAULT 'text',
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_chat_messages_user_id ON chat_messages(user_id);
CREATE INDEX idx_chat_messages_created_at ON chat_messages(created_at);
```

#### 4.2.2 User States (PostgreSQL)
```sql
CREATE TABLE user_states (
    user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    current_state VARCHAR(100) NOT NULL DEFAULT 'start',
    context JSONB DEFAULT '{}',
    platform VARCHAR(20) DEFAULT 'flutter',
    updated_at TIMESTAMP DEFAULT NOW()
);
```

#### 4.2.3 Daily Plans (PostgreSQL)
```sql
CREATE TABLE daily_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    plan_date DATE NOT NULL,
    plan_data JSONB NOT NULL,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, plan_date)
);

CREATE INDEX idx_daily_plans_user_date ON daily_plans(user_id, plan_date);
```

#### 4.2.4 Special Scenarios (PostgreSQL)
```sql
CREATE TABLE special_scenarios (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    scenario_data JSONB NOT NULL,
    selected_option INTEGER DEFAULT 0,
    status VARCHAR(20) DEFAULT 'pending_selection',
    event_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
```

#### 4.2.5 Push Jobs (PostgreSQL)
```sql
CREATE TABLE push_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    job_type VARCHAR(50) NOT NULL,
    message_text TEXT NOT NULL,
    send_at TIMESTAMP NOT NULL,
    status VARCHAR(20) DEFAULT 'scheduled',
    sent_at TIMESTAMP,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_push_jobs_send_at ON push_jobs(send_at) WHERE status = 'scheduled';
```

### 4.3 Структура AI Assistant Service (трансформация Scenario Service)

```
services/ai-assistant-service/  # Rename from scenario-service
├── cmd/
│   └── main.go                 # ♻️  Существующий
├── internal/
│   ├── config/
│   │   └── config.go          # ♻️  Существующий
│   ├── handlers/
│   │   ├── scenario_handler.go # ♻️  Существующий (сохраняем)
│   │   ├── chat_handler.go     # 🆕 Новый (чат API)
│   │   └── assistant_handler.go # 🆕 Новый (анкеты, планы)
│   ├── services/
│   │   ├── scenario_service.go # ♻️  Существующий
│   │   ├── chat_service.go     # 🆕 Новый
│   │   ├── assistant_service.go # 🆕 Новый (основная логика)
│   │   └── state_machine.go    # 🆕 Новый
│   ├── models/
│   │   ├── scenario.go         # ♻️  Существующий
│   │   ├── chat.go            # 🆕 Новый
│   │   ├── plan.go            # 🆕 Новый
│   │   └── user_state.go      # 🆕 Новый
│   └── storage/
│       ├── scenario_storage.go # ♻️  Существующий
│       ├── chat_storage.go     # 🆕 Новый
│       ├── plan_storage.go     # 🆕 Новый
│       └── state_storage.go    # 🆕 Новый
├── migrations/                 # ♻️  Расширяем существующие
└── Dockerfile                  # ♻️  Существующий
```

### 4.4 AI Service (новый микросервис)

```
services/ai-service/
├── cmd/
│   └── main.go
├── internal/
│   ├── config/
│   ├── handlers/
│   │   └── ai_handler.go
│   ├── services/
│   │   ├── gemini_service.go
│   │   └── parser_service.go
│   └── models/
│       ├── plan_request.go
│       └── scenario_request.go
└── Dockerfile
```

### 4.5 Plan Service (новый микросервис)

```
services/plan-service/
├── cmd/
│   └── main.go
├── internal/
│   ├── config/
│   ├── handlers/
│   │   └── plan_handler.go
│   ├── services/
│   │   └── plan_service.go
│   ├── models/
│   │   ├── daily_plan.go
│   │   └── special_scenario.go
│   └── storage/
│       └── plan_storage.go
└── Dockerfile
```

## 5. API Спецификация

### 5.1 AI Assistant Service Endpoints

#### 5.1.1 Chat API
```
POST /chat/message
Body: {
  "user_id": 123,
  "message": "Помоги составить план на завтра"
}
Response: {
  "id": "uuid",
  "response": "Конечно! Начнем с анкеты...",
  "type": "survey",
  "metadata": {
    "buttons": [{"text": "Начать", "action": "start_survey"}]
  }
}

GET /chat/history/{user_id}?limit=50&offset=0
Response: {
  "messages": [...],
  "total": 150
}

POST /chat/clear/{user_id}
Response: {"status": "success"}
```

#### 5.1.2 Assistant API
```
POST /assistant/survey           # Обработка анкет
POST /assistant/scenarios        # Генерация спец-сценариев  
GET  /assistant/plans/{user_id}  # Планы пользователя
```

#### 5.1.3 Scenario API (существующий)
```
POST /search      # ♻️  Существующий
POST /feedback    # ♻️  Существующий
```

### 5.2 AI Service Endpoints

```
POST /ai/generate-plan
Body: {
  "survey_data": {...},
  "user_profile": {...}
}

POST /ai/generate-scenario
Body: {
  "scenario_data": {...},
  "user_profile": {...}
}

POST /ai/parse-time
Body: {
  "text": "завтра в 8 утра",
  "context": {...}
}
```

### 5.3 Plan Service Endpoints

```
POST /plans/daily
GET  /plans/daily/{user_id}?date=2025-01-15
POST /plans/special
GET  /plans/special/{user_id}
PUT  /plans/special/{id}/select-option
```

## 6. Интеграция сервисов

### 6.1 Взаимодействие через API Gateway

```
Flutter Chat → API Gateway → AI Assistant Service → {
  ├── AI Service (для генерации планов)
  ├── Plan Service (для сохранения планов)  
  ├── User Service (для профилей)
  └── Notification Service (для уведомлений)
}
```

### 6.2 Пример потока данных

#### 6.2.1 Создание ежедневного плана
```
1. Flutter → POST /chat/message "Составь план на завтра"
2. AI Assistant → State Machine → "evening_survey_start"
3. AI Assistant → Response "Во сколько проснешься?"
4. Flutter → POST /chat/message "В 8 утра"
5. AI Assistant → AI Service /parse-time → 08:00
6. ... (собираем всю анкету)
7. AI Assistant → AI Service /generate-plan
8. AI Assistant → Plan Service /plans/daily
9. AI Assistant → Notification Service /schedule-jobs
10. Flutter ← Response с планом
```

## 7. Развертывание

### 7.1 Docker Compose обновление

```yaml
# Добавляем к существующему docker-compose.yml
services:
  # Существующие сервисы остаются
  api-gateway: # ✅ Без изменений
  auth-service: # ✅ Без изменений  
  postgres: # ✅ Без изменений
  
  # Переименовываем
  ai-assistant-service: # 🔄 Was scenario-service
    build: ./services/ai-assistant-service
    depends_on: [postgres]
    environment:
      - DB_URL=postgres://postgres:password@postgres:5432/lunary
      - AI_SERVICE_URL=http://ai-service:8088
      - PLAN_SERVICE_URL=http://plan-service:8089
    ports:
      - "8082:8082"
  
  # Новые сервисы
  ai-service:
    build: ./services/ai-service
    environment:
      - GEMINI_API_KEY=${GEMINI_API_KEY}
    ports:
      - "8088:8088"
      
  plan-service:
    build: ./services/plan-service
    depends_on: [postgres]
    environment:
      - DB_URL=postgres://postgres:password@postgres:5432/lunary
    ports:
      - "8089:8089"
```

### 7.2 Миграция данных

```sql
-- 1. Сохраняем существующие данные scenarios, steps, feedback
-- 2. Добавляем новые таблицы
-- 3. Обновляем API Gateway роуты
-- 4. Постепенный rollout новых endpoints
```

## 8. Этапы разработки

### 8.1 Этап 1: Подготовка инфраструктуры (1-2 недели)
- [x] Расширение Auth Service с новыми полями профиля  
- [x] Создание новых таблиц в PostgreSQL
- [x] Настройка миграций

### 8.2 Этап 2: AI Service (2-3 недели) 
- [x] Создание AI Service микросервиса
- [x] Интеграция с Google Gemini
- [x] Парсинг времени и генерация планов
- [x] Unit тесты

### 8.3 Этап 3: Plan Service (2-3 недели)
- [x] Создание Plan Service микросервиса  
- [x] CRUD операции для планов
- [x] Integration с PostgreSQL JSONB
- [x] API endpoints

### 8.4 Этап 4: AI Assistant Service (3-4 недели)
- [x] Трансформация Scenario Service
- [x] Chat API endpoints
- [x] State Machine для диалогов
- [x] Интеграция с AI и Plan сервисами

### 8.5 Этап 5: Notification Service доработка (1-2 недели)
- [x] Планировщик push-уведомлений
- [x] Интеграция с планами
- [x] Система напоминаний

### 8.6 Этап 6: Тестирование и интеграция (2-3 недели)
- [x] Integration тесты между сервисами
- [x] E2E тесты API
- [x] Нагрузочное тестирование
- [x] Documentation

## 9. Мониторинг и метрики

### 9.1 Специфичные метрики для чата
```go
// Metrics для AI Assistant
ChatMessagesTotal       - количество сообщений
ChatResponseTime        - время ответа ассистента  
AIRequestsTotal         - обращения к Gemini
PlanGenerationSuccess   - успешная генерация планов
StateTransitions        - переходы между состояниями
```

### 9.2 Health Checks
```
GET /health - AI Assistant Service
GET /health - AI Service  
GET /health - Plan Service
```

## 10. Безопасность

### 10.1 Аутентификация
- JWT токены от Auth Service для всех запросов
- Проверка user_id в каждом chat message
- Rate limiting на AI запросы

### 10.2 Валидация данных
- Sanitization всех пользовательских сообщений
- Валидация JSON структур планов
- Защита от injection в JSONB запросах

## 11. Производительность

### 11.1 Кеширование
- Redis для частых user states (опционально)
- PostgreSQL connection pooling
- HTTP response caching для статических сценариев

### 11.2 Масштабирование
- Horizontal scaling каждого сервиса независимо
- Database read replicas для chat history
- AI Service rate limiting per user

## 12. Заключение

### 12.1 Преимущества новой архитектуры
✅ **Упрощение**: Нет Telegram Bot сложности  
✅ **Интеграция**: Нативный опыт в приложении
✅ **Переиспользование**: 70% существующей инфраструктуры
✅ **Масштабируемость**: Микросервисная архитектура
✅ **Скорость разработки**: 8-12 недель вместо 20-24

### 12.2 Итоговые временные рамки
- **MVP чата**: 6-8 недель
- **Полный функционал**: 10-12 недель  
- **Production ready**: 12-14 недель

### 12.3 Команда
- **Senior Go Developer** (1 чел.) - основная разработка
- **AI Engineer** (0.5 чел.) - Gemini интеграция  
- **QA Engineer** (0.5 чел.) - тестирование
- **DevOps** (0.3 чел.) - поддержка инфраструктуры

---

**Версия документа:** 2.0 (Обновленная архитектура)  
**Дата обновления:** 2025-01-15  
**Основные изменения:** Замена Telegram Bot на встроенный AI-чат, использование PostgreSQL вместо MongoDB
