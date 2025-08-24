# User Service

Микросервис для управления профилями пользователей в платформе Lunary AI Assistant.

## 📋 Описание

User Service отвечает за:
- Создание и управление профилями пользователей
- Хранение персональных данных (имя, профессия, город и т.д.)
- Настройки пользователей (язык, часовой пояс, тон уведомлений)
- Интеграция с auth-service для проверки JWT токенов

## 🏗️ Архитектура

User Service следует той же архитектуре что и auth-service:

```
user-service/
├── cmd/
│   └── main.go                 # Точка входа + graceful shutdown
├── internal/
│   ├── handler/
│   │   └── handler.go          # HTTP обработчики + middleware
│   ├── model/
│   │   └── model.go            # Модели данных
│   ├── service/
│   │   └── service.go          # Бизнес-логика + валидация
│   ├── storage/
│   │   └── storage.go          # Работа с БД + конфигурация
│   └── token/
│       └── token.go            # JWT обработка + контекст
├── migrations/
│   ├── 001_create_user_profiles_table.up.sql
│   └── 001_create_user_profiles_table.down.sql
├── Dockerfile
├── go.mod
└── README.md
```

### Ключевые компоненты:

- **Storage Layer**: Чистый SQL (database/sql) вместо ORM
- **Service Layer**: Валидация через go-playground/validator
- **Handler Layer**: Чистый net/http + gorilla/mux
- **Common Integration**: Использует shared логирование, ошибки и health checks
- **JWT Integration**: Совместимость с auth-service токенами

## 🚀 API Endpoints

### Health Check
```http
GET /health
GET /health/live
GET /health/ready
```

**Пример ответа:**
```json
{
  "status": "ok",
  "service": "user-service", 
  "version": "1.0.0",
  "checks": {
    "database": "ok"
  }
}
```

### Профили пользователей (требует авторизации)

#### Создать профиль
```http
POST /api/v1/profile
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json

{
  "name": "Иван",
  "lastname": "Иванов", 
  "sex": "MALE",
  "profession": "Developer",
  "city": "Almaty",
  "language": "RU",
  "notification_tone": "NEUTRAL"
}
```

**Ответ (201 Created):**
```json
{
  "message": "Профиль создан успешно"
}
```

#### Получить профиль
```http
GET /api/v1/profile
Authorization: Bearer <JWT_TOKEN>
```

**Ответ (200 OK):**
```json
{
  "message": "Профиль получен успешно",
  "profile": {
    "id": 123,
    "phone": "+77771234567",
    "email": "user@example.com",
    "name": "Иван",
    "lastname": "Иванов",
    "sex": "MALE",
    "profession": "Developer",
    "city": "Almaty",
    "timezone": "Asia/Almaty",
    "language": "RU",
    "notification_tone": "NEUTRAL",
    "last_special_scenario_date": null,
    "created_at": "2025-01-15T10:00:00Z"
  }
}
```

#### Обновить профиль
```http
PUT /api/v1/profile
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json

{
  "name": "Новое имя",
  "city": "Astana"
}
```

**Ответ (200 OK):**
```json
{
  "message": "Профиль обновлен успешно"
}
```

#### Удалить профиль
```http
DELETE /api/v1/profile
Authorization: Bearer <JWT_TOKEN>
```

**Ответ (200 OK):**
```json
{
  "message": "Профиль удален успешно"
}
```

### Обработка ошибок

Все ошибки возвращаются в стандартизированном формате:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Ошибка валидации: Name превышает максимальную длину",
    "type": "validation"
  }
}
```

**Типы ошибок:**
- `validation` (400) - Ошибки валидации данных
- `authentication` (401) - Проблемы с авторизацией  
- `not_found` (404) - Ресурс не найден
- `conflict` (409) - Конфликт данных
- `internal` (500) - Внутренние ошибки сервера

## 🗄️ База данных

### Таблица user_profiles

| Поле | Тип | Описание | Ограничения |
|------|-----|----------|-------------|
| id | SERIAL | Primary Key | NOT NULL |
| user_id | INTEGER | Foreign Key на users.id | UNIQUE, CASCADE DELETE |
| name | VARCHAR(50) | Имя пользователя | - |
| lastname | VARCHAR(50) | Фамилия | - |
| sex | VARCHAR(10) | Пол | CHECK: MALE/FEMALE |
| profession | VARCHAR(100) | Профессия | - |
| city | VARCHAR(100) | Город | DEFAULT: 'Almaty' |
| timezone | VARCHAR(50) | Часовой пояс | DEFAULT: 'Asia/Almaty' |
| language | VARCHAR(5) | Язык | CHECK: RU/KZ/EN, DEFAULT: 'RU' |
| notification_tone | VARCHAR(20) | Тон уведомлений | CHECK: SOFT/NEUTRAL/STRICT, DEFAULT: 'NEUTRAL' |
| last_special_scenario_date | TIMESTAMP | Дата последнего специального сценария | - |
| created_at | TIMESTAMP | Дата создания | DEFAULT: CURRENT_TIMESTAMP |
| updated_at | TIMESTAMP | Дата обновления | AUTO UPDATE |

### Связи с другими таблицами:

- **users.id** → **user_profiles.user_id** (1:1, CASCADE DELETE)

### Индексы:

- `idx_user_profiles_user_id` (UNIQUE) - связь с users
- `idx_user_profiles_city` - поиск по городам
- `idx_user_profiles_language` - поиск по языкам

## ⚙️ Конфигурация

### Переменные окружения

| Переменная | По умолчанию | Описание |
|------------|--------------|----------|
| PORT | 8083 | Порт сервиса |
| DB_HOST | localhost | Хост PostgreSQL |
| DB_PORT | 5432 | Порт PostgreSQL |
| DB_USER | postgres | Пользователь БД |
| DB_PASSWORD | postgres | Пароль БД |
| DB_NAME | lunary_db | Имя базы данных |
| DB_SSL_MODE | disable | SSL режим для БД |
| JWT_SECRET | your-secret-key | Секретный ключ для JWT (должен совпадать с auth-service) |
| AUTH_SERVICE_URL | http://localhost:8081 | URL auth-service |

### Пример .env файла
```bash
# Сервис
PORT=8083

# База данных
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=lunary_db
DB_SSL_MODE=disable

# JWT (ВАЖНО: должен совпадать с auth-service)
JWT_SECRET=your-super-secret-key-here

# Интеграция
AUTH_SERVICE_URL=http://localhost:8081
```

## 🛠️ Разработка

### Требования
- Go 1.24.2+
- PostgreSQL 15+
- Auth Service (для JWT токенов)

### Локальный запуск

1. **Клонирование и переход:**
```bash
cd services/user-service
```

2. **Установка зависимостей:**
```bash
go mod tidy
```

3. **Настройка базы данных:**
```sql
-- Подключиться к PostgreSQL
psql -U postgres

-- Создать базу данных (если еще нет)
CREATE DATABASE lunary_db;

-- Переключиться на базу
\c lunary_db

-- Запустить миграции
\i migrations/001_create_user_profiles_table.up.sql
```

4. **Настройка переменных окружения:**
```bash
# Создать .env файл или экспортировать переменные
export JWT_SECRET="same-secret-as-auth-service"
export DB_PASSWORD="your-db-password"
```

5. **Запуск сервиса:**
```bash
go run cmd/main.go
```

Сервис будет доступен на `http://localhost:8083`

### Тестирование

```bash
# Health check
curl http://localhost:8083/health

# Получение JWT токена из auth-service
JWT_TOKEN=$(curl -s -X POST http://localhost:8081/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"phone":"+77001234567","password":"password"}' \
  | jq -r '.token')

# Создание профиля
curl -X POST http://localhost:8083/api/v1/profile \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","city":"Almaty","language":"RU"}'

# Получение профиля
curl http://localhost:8083/api/v1/profile \
  -H "Authorization: Bearer $JWT_TOKEN"
```

## 🐳 Docker

### Сборка образа
```bash
docker build -t user-service .
```

### Запуск контейнера
```bash
docker run -p 8083:8083 \
  -e DB_HOST=host.docker.internal \
  -e DB_PASSWORD=postgres \
  -e JWT_SECRET=your-secret \
  user-service
```

### Docker Compose

```yaml
version: '3.8'
services:
  user-service:
    build: .
    ports:
      - "8083:8083"
    environment:
      - DB_HOST=postgres
      - DB_PASSWORD=postgres
      - JWT_SECRET=your-secret-key
    depends_on:
      - postgres
      - auth-service
    networks:
      - lunary-network

  postgres:
    image: postgres:15
    environment:
      - POSTGRES_DB=lunary_db
      - POSTGRES_PASSWORD=postgres
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./migrations:/docker-entrypoint-initdb.d
    networks:
      - lunary-network

  auth-service:
    build: ../auth-service
    ports:
      - "8081:8081"
    environment:
      - DB_HOST=postgres
      - JWT_SECRET=your-secret-key
    depends_on:
      - postgres
    networks:
      - lunary-network

volumes:
  postgres_data:

networks:
  lunary-network:
    driver: bridge
```

## 🔒 Безопасность

### JWT Авторизация
- Все API endpoints (кроме health checks) требуют JWT токена
- Токены должны быть получены от auth-service
- Токены проверяются через middleware на каждый запрос
- user_id извлекается из токена и используется для операций

### Валидация данных
- Все входные данные валидируются через go-playground/validator
- Максимальные длины полей соответствуют схеме БД
- Проверка допустимых значений для enum полей (sex, language, notification_tone)

### База данных
- Использование prepared statements для защиты от SQL injection
- Foreign Key constraints для целостности данных
- Connection pooling с ограничениями
- Graceful shutdown для завершения транзакций

## 📊 Мониторинг и логирование

### Health Checks
- `/health` - общий статус + проверка БД
- `/health/live` - проверка готовности (liveness probe)
- `/health/ready` - проверка готовности к обслуживанию (readiness probe)

### Логирование
Структурированные логи с контекстом:
```json
{
  "level": "info",
  "message": "User profile created successfully",
  "user_id": 123,
  "profile_id": 456,
  "timestamp": "2025-01-15T10:00:00Z"
}
```

**Логируемые события:**
- Запуск/остановка сервиса
- Database connections/errors
- Profile operations (create/update/delete)
- Authentication events
- Validation errors

### Метрики
Готов к интеграции с Prometheus через:
- HTTP response times
- Request counts по endpoints
- Database connection pool metrics
- Error rates по типам

## 🔄 Интеграция

### С Auth Service
- **JWT токены**: Использует те же ключи что и auth-service
- **User validation**: Проверяет существование пользователя через БД
- **Claims extraction**: Извлекает user_id из JWT для операций

### С Frontend/Mobile
- **REST API**: Стандартизированные JSON endpoints
- **CORS**: Настроен для cross-origin запросов
- **Error handling**: Консистентные коды ошибок
- **Request/Response**: Predictable JSON структуры

### С Database
- **Connection pooling**: Оптимизированные настройки подключений
- **Transactions**: Atomic операции для data integrity
- **Migrations**: Versioned schema changes
- **Indexing**: Оптимизированные запросы

## 📈 Производительность

### Database Optimization
- **Connection Pool**: Max 25 open, 5 idle connections
- **Connection Lifetime**: 5 minutes rotation
- **Prepared Statements**: Повторное использование запросов
- **Indexes**: Оптимизированы для частых запросов

### HTTP Optimization  
- **Timeouts**: Read 15s, Write 15s, Idle 60s
- **Request Limits**: 1MB max request size
- **Keep-Alive**: Persistent connections
- **Graceful Shutdown**: 30s timeout для завершения запросов

## 🚀 Deployment

### Production Checklist
- [ ] Настроены environment variables
- [ ] Database migrations применены
- [ ] JWT_SECRET синхронизирован с auth-service
- [ ] Health checks настроены в load balancer
- [ ] Logging forwarding настроен
- [ ] Metrics collection настроен
- [ ] Resource limits установлены
- [ ] Backup strategy настроена

### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: user-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: user-service
  template:
    metadata:
      labels:
        app: user-service
    spec:
      containers:
      - name: user-service
        image: user-service:latest
        ports:
        - containerPort: 8083
        env:
        - name: DB_HOST
          value: "postgres-service"
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: jwt-secret
              key: secret
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8083
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8083
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "64Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "500m"
```

## 📝 Changelog

### v1.0.0 (2025-01-15)
- ✅ Базовая реализация CRUD операций для профилей
- ✅ JWT авторизация совместимая с auth-service
- ✅ PostgreSQL интеграция с чистым SQL
- ✅ Валидация через go-playground/validator
- ✅ Структурированное логирование
- ✅ Health checks (health/live/ready)
- ✅ Graceful shutdown с таймаутами
- ✅ Docker поддержка
- ✅ CORS middleware для frontend
- ✅ Типизированные ошибки
- ✅ Connection pooling оптимизация

## 🤝 Contributing

1. Создайте feature branch от main
2. Внесите изменения с соблюдением архитектурных принципов
3. Добавьте/обновите тесты (если необходимо)
4. Убедитесь что логирование и error handling добавлены
5. Обновите документацию при изменении API
6. Создайте Pull Request с описанием изменений

### Code Style
- Следуйте Go conventions (gofmt, golint)
- Используйте структурированное логирование через shared
- Все ошибки должны быть типизированными (shared.AppError)
- Валидация обязательна для всех входных данных
- Database операции должны использовать prepared statements

## 📄 License

MIT License - см. файл LICENSE для деталей.