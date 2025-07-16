# LUNARY Backend Services

## Обзор

LUNARY - это платформа для планирования мероприятий с использованием ИИ, построенная на архитектуре микросервисов. Бэкенд состоит из нескольких независимых сервисов, которые взаимодействуют друг с другом через REST API и gRPC.

## Архитектура микросервисов

```
┌─────────────────┐    ┌──────────────────┐    ┌────────────────┐
│   Frontend      │───▶│   API Gateway    │───▶│   Auth Service │
│   (Flutter)     │    │   (Port: 8080)   │    │   (Port: 8081) │
└─────────────────┘    └──────────────────┘    └────────────────┘
                              │
                              ├─────────────────────────────────┐
                              │                                 │
                    ┌─────────▼──────────┐          ┌─────────▼──────────┐
                    │  Scenario Service  │          │   User Service     │
                    │   (Port: 8082)     │          │   (Port: 8083)     │
                    └────────────────────┘          └────────────────────┘
                              │                                 │
                    ┌─────────▼──────────┐          ┌─────────▼──────────┐
                    │   Ad Service       │          │  Media Service     │
                    │   (Port: 8084)     │          │   (Port: 8085)     │
                    └────────────────────┘          └────────────────────┘
                              │
                    ┌─────────▼──────────┐
                    │ Notification Service│
                    │   (Port: 8086)     │
                    └────────────────────┘
                              │
                    ┌─────────▼──────────┐
                    │    PostgreSQL      │
                    │   (Port: 5432)     │
                    └────────────────────┘
```

## Сервисы

### Основные сервисы (MVP)

| Сервис | Порт | Статус | Описание |
|--------|------|--------|----------|
| **API Gateway** | 8080 | ✅ Готов | Маршрутизация запросов, аутентификация, rate limiting |
| **Auth Service** | 8081 | ✅ Готов | Регистрация, авторизация, JWT токены |
| **Scenario Service** | 8082 | ✅ Готов | Поиск и управление сценариями мероприятий |

### Планируемые сервисы

| Сервис | Порт | Статус | Описание |
|--------|------|--------|----------|
| **User Service** | 8083 | 🔄 В разработке | Управление профилями пользователей |
| **Ad Service** | 8084 | 🔄 В разработке | Управление рекламными объявлениями |
| **Media Service** | 8085 | 🔄 В разработке | Загрузка и хранение медиафайлов |
| **Notification Service** | 8086 | 🔄 В разработке | Push-уведомления и email рассылки |

## Быстрый старт

### Предварительные требования

- Docker и Docker Compose
- Go 1.21+ (для разработки)
- PostgreSQL (через Docker)

### Запуск всех сервисов

1. **Клонировать репозиторий**
   ```bash
   git clone <repository-url>
   cd LUNARY
   ```

2. **Запустить через Docker Compose**
   ```bash
   docker-compose up --build
   ```

3. **Проверить статус сервисов**
   ```bash
   # API Gateway
   curl http://localhost:8080/health
   
   # Auth Service
   curl http://localhost:8081/health
   
   # Scenario Service
   curl http://localhost:8082/health
   ```

### Переменные окружения

Каждый сервис использует собственный `.env` файл:
- `api-gateway/.env`
- `auth-service/.env`
- `scenario-service/.env`
- и т.д.

Основные переменные:
```env
DB_HOST=postgres
DB_PORT=5432
DB_USER=lunary_user
DB_PASSWORD=password
DB_NAME=lunary_db
JWT_SECRET=your-secret-key
```

## База данных

### PostgreSQL

- **Хост**: localhost:5432 (или postgres:5432 в Docker)
- **База данных**: lunary_db
- **Пользователь**: lunary_user
- **Пароль**: password

### Миграции

Миграции запускаются автоматически при старте сервисов:
- Auth Service: создание таблицы users
- Scenario Service: создание таблиц scenarios, feedback + начальные данные

## Разработка

### Структура проекта

```
services/
├── api-gateway/           # API Gateway (Go + Gin)
├── auth-service/          # Сервис аутентификации (Go + Gin)
├── scenario-service/      # Сервис сценариев (Go + Gin)
├── user-service/          # Сервис пользователей (планируется)
├── ad-service/           # Сервис рекламы (планируется)
├── media-service/        # Сервис медиа (планируется)
├── notification-service/ # Сервис уведомлений (планируется)
└── common/              # Общие компоненты (протоколы, утилиты)
```

### Запуск отдельного сервиса

```bash
cd services/auth-service
go run cmd/main.go
```

### Тестирование

```bash
# Запуск тестов для всех сервисов
cd services
for service in */; do
  if [ -f "$service/go.mod" ]; then
    echo "Testing $service"
    cd "$service" && go test ./... && cd ..
  fi
done
```

## API Documentation

### Swagger UI

- **Auth Service**: http://localhost:8081/swagger/index.html
- **Scenario Service**: http://localhost:8082/swagger/index.html (планируется)

### Postman Collection

Импортируйте коллекцию из `../TESTING.md` для полного набора API тестов.

## Мониторинг и логирование

### Логи

```bash
# Просмотр логов всех сервисов
docker-compose logs -f

# Просмотр логов конкретного сервиса
docker-compose logs -f auth-service
```

### Health Checks

Все сервисы предоставляют endpoint `/health`:
- API Gateway: http://localhost:8080/health
- Auth Service: http://localhost:8081/health
- Scenario Service: http://localhost:8082/health

## Развертывание

### Production

1. **Настройка переменных окружения**
   - Используйте безопасные пароли
   - Настройте внешний PostgreSQL
   - Настройте SSL сертификаты

2. **Docker Compose Production**
   ```bash
   docker-compose -f docker-compose.prod.yml up -d
   ```

3. **Kubernetes** (в планах)
   ```bash
   kubectl apply -f deploy/k8s/
   ```

## Безопасность

- JWT токены для аутентификации
- Rate limiting в API Gateway
- Валидация входных данных
- Хеширование паролей (bcrypt)
- CORS настройки

## Производительность

- Пулы соединений к БД
- HTTP/2 поддержка
- Кеширование (планируется Redis)
- Horizontal scaling готовность

## Связанные документы

- [API Gateway README](./api-gateway/README.md)
- [Auth Service README](./auth-service/README.md)
- [Scenario Service README](./scenario-service/README.md)
- [User Service README](./user-service/README.md)
- [Ad Service README](./ad-service/README.md)
- [Media Service README](./media-service/README.md)
- [Notification Service README](./notification-service/README.md)
- [Тестирование API](../TESTING.md)

## Поддержка

Для вопросов и предложений создавайте issues в репозитории проекта.

## Лицензия

См. файл [LICENSE](../LICENSE) в корне проекта. 