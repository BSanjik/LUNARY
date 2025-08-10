# 🌙 LUNARY

**Современная микросервисная платформа для объявлений с AI-рекомендациями**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Flutter](https://img.shields.io/badge/Flutter-3.0+-02569B?style=flat&logo=flutter)](https://flutter.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-316192?style=flat&logo=postgresql)](https://postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](https://docker.com/)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-Ready-326CE5?style=flat&logo=kubernetes)](https://kubernetes.io/)

---

## 📋 Описание проекта

**LUNARY** — это production-ready микросервисная платформа, которая позволяет пользователям публиковать объявления, получать AI-рекомендации и управлять медиафайлами. Проект построен с использованием современных технологий и лучших практик enterprise-разработки.

### 🎯 Ключевые особенности

- ✅ **Микросервисная архитектура** с независимым масштабированием
- ✅ **Production-ready** код с полным покрытием тестами
- ✅ **Структурированное логирование** в JSON формате
- ✅ **Health checks** для мониторинга состояния сервисов
- ✅ **Graceful shutdown** для корректного завершения
- ✅ **API документация** с Swagger/OpenAPI
- ✅ **Контейнеризация** с Docker и Kubernetes
- ✅ **Валидация запросов** с детальными ошибками
- ✅ **JWT аутентификация** с refresh токенами

---

## 🏗️ Архитектура

### Микросервисы (`services/`)

| Сервис | Порт | Описание | Статус |
|--------|------|----------|---------|
| **api-gateway** | 8080 | Единая точка входа, маршрутизация | ✅ Ready |
| **auth-service** | 8081 | Регистрация, авторизация, JWT | ✅ Production |
| **scenario-service** | 8082 | Управление сценариями | ✅ Ready |
| **user-service** | 8083 | Профили пользователей | 🚧 In Progress |
| **ad-service** | 8084 | Создание и поиск объявлений | 🚧 In Progress |
| **media-service** | 8085 | Загрузка медиафайлов | 🚧 In Progress |
| **notification-service** | 8086 | Email и push уведомления | 🚧 In Progress |

### Фронтенд (`frontend/`)
- **React Native** мобильное приложение
- **TypeScript** для типобезопасности
- **Redux Toolkit** для управления состоянием
- **React Navigation** для навигации
- **Адаптивный дизайн** для iOS и Android

### Инфраструктура (`deploy/`)
- **Docker Compose** для локальной разработки
- **Kubernetes** манифесты для продакшена
- **Health checks** и мониторинг

---

## 🛠️ Технологии

### Backend
- **Go 1.21+** - основной язык разработки
- **Gorilla/mux** - HTTP роутинг
- **PostgreSQL 15** - основная база данных
- **JWT** - аутентификация и авторизация
- **Logrus** - структурированное логирование
- **GORM** - ORM для работы с БД
- **go-playground/validator** - валидация запросов

### Frontend
- **React Native 0.72+** - мобильное приложение
- **TypeScript** - типобезопасность
- **Redux Toolkit** - управление состоянием  
- **React Navigation 6** - навигация
- **React Hook Form** - управление формами
- **Axios** - HTTP клиент для API

### DevOps
- **Docker** - контейнеризация
- **Kubernetes** - оркестрация
- **Docker Compose** - локальная разработка
- **GitHub Actions** - CI/CD (планируется)

### Monitoring & Observability
- **Health Checks** - проверка состояния сервисов
- **Structured Logging** - JSON логи для мониторинга
- **Graceful Shutdown** - корректное завершение
- **Prometheus** готовность (планируется)

---

## 🚀 Быстрый старт

### Требования

- **Go 1.21+**
- **Node.js 16+** 
- **React Native CLI**
- **Docker & Docker Compose**
- **PostgreSQL 15** (или Docker)
- **Android Studio** (для Android разработки)
- **Xcode** (для iOS разработки, только macOS)

### 1. Клонирование репозитория

```bash
git clone https://github.com/BSanjik/LUNARY.git
cd LUNARY
```

### 2. Запуск с Docker Compose (рекомендуется)

```bash
# Создаем .env файл
cp .env.example .env

# Запускаем все сервисы
docker-compose up -d

# Проверяем статус
docker-compose ps
```

### 3. Локальная разработка

```bash
# 1. Запуск PostgreSQL
docker run -d --name postgres-lunary \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=lunary \
  -p 5432:5432 \
  postgres:15

# 2. Настройка auth-service
cd services/auth-service
cp .env.example .env
go mod download
go run cmd/main.go

# 3. Запуск React Native приложения
cd frontend
npm install

# Для iOS (только macOS)
cd ios && pod install && cd ..

# Запуск Metro bundler
npm start

# В другом терминале - запуск на устройстве
npm run android  # Для Android
npm run ios      # Для iOS (только macOS)
```

### 4. Проверка работоспособности

```bash
# Health checks
curl http://localhost:8081/health
curl http://localhost:8081/health/ready

# API документация
open http://localhost:8081/swagger/

# Тестовый запрос
curl -X POST http://localhost:8081/api/v1/registration \
  -H "Content-Type: application/json" \
  -d '{"phone": "+77001234567", "password": "password123"}'
```

---

## 📊 API Endpoints

### Auth Service (8081)

| Method | Endpoint | Описание |
|--------|----------|----------|
| POST | `/api/v1/registration` | Регистрация пользователя |
| POST | `/api/v1/login` | Авторизация пользователя |
| GET | `/health` | Детальное состояние сервиса |
| GET | `/health/live` | Проверка "живости" |
| GET | `/health/ready` | Готовность к работе |
| GET | `/swagger/` | API документация |

### Scenario Service (8082)

| Method | Endpoint | Описание |
|--------|----------|----------|
| GET | `/api/v1/scenarios` | Список сценариев |
| GET | `/api/v1/scenarios/{id}` | Получение сценария |
| POST | `/api/v1/scenarios/{id}/feedback` | Отправка отзыва |

---

## 🧪 Тестирование

### Unit тесты

```bash
# Все тесты
go test ./...

# Конкретный сервис
cd services/auth-service
go test ./internal/...

# С покрытием
go test -cover ./...
```

### Health Checks

```bash
# Тест health checks
cd services/auth-service
bash test_health.sh
```

### Graceful Shutdown

```bash
# Тест graceful shutdown
cd services/auth-service
bash test_graceful_shutdown.sh
```

### Postman коллекция

```bash
# Импорт коллекции
newman run docs/postman/lunary-collection.json
```

---

## 📁 Структура проекта

```
LUNARY/
├── 📁 deploy/                 # Конфигурации деплоя
│   ├── docker-compose.yml    # Docker Compose
│   └── k8s/                  # Kubernetes манифесты
├── 📁 docs/                  # Документация
│   ├── api-spec.md          # API спецификация
│   ├── architecture.md      # Архитектура системы
│   ├── graceful-shutdown.md  # Graceful shutdown
│   └── tech-stack.md        # Технологический стек
├── 📁 frontend/              # React Native приложение
│   ├── src/                 # Исходный код TypeScript
│   │   ├── components/      # Переиспользуемые компоненты
│   │   ├── screens/         # Экраны приложения
│   │   ├── navigation/      # Навигация
│   │   ├── services/        # API сервисы
│   │   ├── store/          # Redux store
│   │   └── types/          # TypeScript типы
│   ├── android/             # Android платформа
│   ├── ios/                 # iOS платформа
│   └── package.json         # Зависимости Node.js
├── 📁 services/              # Go микросервисы
│   ├── common/              # Общие компоненты
│   │   └── shared/          # Логирование, ошибки, health checks
│   ├── auth-service/        # ✅ Готов к продакшену
│   ├── scenario-service/    # ✅ Готов к продакшену
│   ├── api-gateway/         # 🚧 В разработке
│   └── .../                 # Остальные сервисы
├── 📄 init.sql              # Схема базы данных
├── 📄 docker-compose.yml    # Локальная разработка
└── 📄 README.md            # Этот файл
```

---

## 🔧 Конфигурация

### Переменные окружения

#### Auth Service
```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=lunary

# JWT
JWT_SECRET=your_super_secret_key_change_in_production

# Server
PORT=8081
LOG_LEVEL=info
```

### Docker Compose

```yaml
# Health checks включены
healthcheck:
  test: ["CMD", "curl", "-f", "http://localhost:8081/health/ready"]
  interval: 30s
  timeout: 10s
  retries: 3
  start_period: 40s
```

### Kubernetes

```yaml
# Graceful shutdown
terminationGracePeriodSeconds: 45

# Health checks
livenessProbe:
  httpGet:
    path: /health/live
    port: 8081
readinessProbe:
  httpGet:
    path: /health/ready
    port: 8081
```

---

## 📊 Мониторинг и логирование

### Структурированные логи

```json
{
  "timestamp": "2024-01-15T10:30:00.000Z",
  "level": "info",
  "message": "User registered successfully",
  "service": "auth-service",
  "user_id": 123,
  "phone": "+77001234567",
  "event": "register",
  "success": true
}
```

### Health Check ответы

```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00.000Z",
  "service": {
    "name": "auth-service",
    "version": "1.0.0",
    "uptime": "2h35m12s"
  },
  "checks": {
    "database": {
      "status": "healthy",
      "duration_ms": 5
    }
  },
  "system_info": {
    "go_version": "go1.21.5",
    "mem_alloc_mb": 12
  }
}
```

---

## 🔄 Git Workflow

### Работа с ветками

```bash
# 1. Обновляем main
git checkout main
git pull origin main

# 2. Создаем ветку для задачи
git checkout -b feature/task-name

# 3. Делаем изменения
git add .
git commit -m "feat: добавить новую функциональность"

# 4. Отправляем на сервер
git push -u origin feature/task-name

# 5. Создаем Pull Request на GitHub
```

### Соглашения о коммитах

- `feat:` - новая функциональность
- `fix:` - исправление багов
- `docs:` - обновление документации
- `test:` - добавление тестов
- `refactor:` - рефакторинг кода

---

## 🎯 Production Ready Features

### ✅ Реализовано

- **Unit Tests** - полное покрытие бизнес-логики
- **Error Handling** - структурированные ошибки с кодами
- **Request Validation** - валидация с детальными сообщениями
- **Structured Logging** - JSON логи с метаданными
- **Health Checks** - liveness, readiness, detailed monitoring
- **Graceful Shutdown** - корректное завершение сервисов
- **API Documentation** - Swagger спецификация
- **Containerization** - Docker images для всех сервисов
- **Database Migrations** - версионирование схемы БД

### 🚧 В планах

- **Rate Limiting** - ограничение количества запросов
- **Circuit Breaker** - защита от каскадных сбоев
- **Distributed Tracing** - трассировка запросов
- **Metrics Collection** - сбор метрик для Prometheus
- **CI/CD Pipeline** - автоматическое тестирование и деплой

---

## 🤝 Контрибьюция

1. **Fork** репозитория
2. Создайте **feature branch** (`git checkout -b feature/amazing-feature`)
3. **Commit** изменения (`git commit -m 'feat: add amazing feature'`)
4. **Push** в ветку (`git push origin feature/amazing-feature`)
5. Создайте **Pull Request**

### Требования к коду

- ✅ Unit тесты для новой функциональности
- ✅ Структурированное логирование
- ✅ Обработка ошибок
- ✅ Валидация входных данных
- ✅ API документация

---

## 📄 Лицензия

Этот проект лицензирован под **MIT License** - подробности в файле [LICENSE](LICENSE).

---

## 📞 Поддержка

- **Issues**: [GitHub Issues](https://github.com/BSanjik/LUNARY/issues)
- **Discussions**: [GitHub Discussions](https://github.com/BSanjik/LUNARY/discussions)
- **Email**: support@lunary.app

---

## ⭐ Благодарности

Спасибо всем контрибьюторам, которые помогают делать LUNARY лучше!

---

<div align="center">

**🌙 Сделано с ❤️ командой LUNARY**

[⭐ Star this repo](https://github.com/BSanjik/LUNARY) • [🐛 Report Bug](https://github.com/BSanjik/LUNARY/issues) • [💡 Request Feature](https://github.com/BSanjik/LUNARY/issues)

</div>