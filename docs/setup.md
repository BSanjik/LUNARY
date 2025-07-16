# Инструкция по установке и настройке LUNARY

## Содержание
- [Системные требования](#системные-требования)
- [Быстрый старт](#быстрый-старт)
- [Детальная установка](#детальная-установка)
- [Настройка разработки](#настройка-разработки)
- [Решение проблем](#решение-проблем)
- [Дополнительные инструменты](#дополнительные-инструменты)

## Системные требования

### Минимальные требования
- **OS**: Windows 10/11, macOS 10.15+, Ubuntu 18.04+
- **RAM**: 4 GB (рекомендуется 8 GB)
- **Диск**: 10 GB свободного места
- **Интернет**: для загрузки зависимостей

### Обязательное ПО

#### 1. Docker & Docker Compose
```bash
# Windows/macOS - Docker Desktop
# Скачать с: https://www.docker.com/products/docker-desktop

# Ubuntu
sudo apt update
sudo apt install docker.io docker-compose
sudo usermod -aG docker $USER
# Перелогиниться после добавления в группу docker
```

#### 2. Git
```bash
# Windows - Git for Windows
# Скачать с: https://git-scm.com/download/win

# macOS
brew install git

# Ubuntu
sudo apt install git
```

### Опциональное ПО (для разработки)

#### 3. Go 1.21+
```bash
# Скачать с: https://golang.org/dl/

# Проверить установку
go version
```

#### 4. Flutter SDK
```bash
# Скачать с: https://flutter.dev/docs/get-started/install

# Проверить установку
flutter doctor
```

#### 5. PostgreSQL клиент (опционально)
```bash
# Windows - pgAdmin
# macOS
brew install postgresql

# Ubuntu
sudo apt install postgresql-client
```

## Быстрый старт

### 1. Клонирование репозитория
```bash
git clone https://github.com/BSanjik/LUNARY.git
cd LUNARY
```

### 2. Запуск всех сервисов
```bash
# Запуск через Docker Compose
docker-compose up --build

# Или в фоновом режиме
docker-compose up --build -d
```

### 3. Проверка работоспособности
```bash
# API Gateway
curl http://localhost:8080/health

# Auth Service
curl http://localhost:8081/health

# Scenario Service
curl http://localhost:8082/health
```

### 4. Открыть документацию
- **Swagger Auth Service**: http://localhost:8081/swagger/index.html
- **API Gateway Health**: http://localhost:8080/health

**🎉 Готово! Все сервисы запущены и готовы к использованию.**

## Детальная установка

### Структура проекта
```
LUNARY/
├── services/              # Бэкенд микросервисы
│   ├── api-gateway/       # API Gateway (8080)
│   ├── auth-service/      # Аутентификация (8081)
│   ├── scenario-service/  # Сценарии (8082)
│   └── ...               # Другие сервисы
├── frontend/             # Flutter приложение
├── docs/                 # Документация
├── docker-compose.yml    # Конфигурация контейнеров
└── .env файлы           # Конфигурация окружения
```

### Настройка переменных окружения

#### 1. API Gateway (.env)
```bash
# services/api-gateway/.env
PORT=8080
AUTH_SERVICE_URL=http://auth-service:8081
SCENARIO_SERVICE_URL=http://scenario-service:8082
USER_SERVICE_URL=http://user-service:8083
AD_SERVICE_URL=http://ad-service:8084
MEDIA_SERVICE_URL=http://media-service:8085
NOTIFICATION_SERVICE_URL=http://notification-service:8086
JWT_SECRET=your-super-secret-jwt-key-change-in-production
RATE_LIMIT=100
```

#### 2. Auth Service (.env)
```bash
# services/auth-service/.env
PORT=8081
DB_HOST=postgres
DB_PORT=5432
DB_USER=lunary_user
DB_PASSWORD=password
DB_NAME=lunary_db
JWT_SECRET=your-super-secret-jwt-key-change-in-production
BCRYPT_COST=12
```

#### 3. Scenario Service (.env)
```bash
# services/scenario-service/.env
PORT=8082
DB_HOST=postgres
DB_PORT=5432
DB_USER=lunary_user
DB_PASSWORD=password
DB_NAME=lunary_db
```

#### 4. PostgreSQL (docker-compose.yml)
```yaml
# Уже настроена в docker-compose.yml
POSTGRES_DB=lunary_db
POSTGRES_USER=lunary_user
POSTGRES_PASSWORD=password
```

### Пошаговый запуск

#### 1. Подготовка окружения
```bash
# Создать необходимые директории
mkdir -p logs
mkdir -p data/postgres

# Проверить Docker
docker --version
docker-compose --version
```

#### 2. Сборка образов
```bash
# Сборка всех сервисов
docker-compose build

# Или по отдельности
docker-compose build api-gateway
docker-compose build auth-service
docker-compose build scenario-service
```

#### 3. Запуск базы данных
```bash
# Запустить только PostgreSQL
docker-compose up postgres -d

# Проверить подключение
docker-compose exec postgres psql -U lunary_user -d lunary_db -c "\dt"
```

#### 4. Запуск сервисов по очереди
```bash
# Auth Service
docker-compose up auth-service -d

# Scenario Service
docker-compose up scenario-service -d

# API Gateway
docker-compose up api-gateway -d

# Проверить статус
docker-compose ps
```

### Миграции базы данных

Миграции запускаются автоматически при старте сервисов. Если нужно запустить вручную:

```bash
# Auth Service миграции
docker-compose exec auth-service ./wait-for-it.sh postgres:5432 -- go run cmd/main.go

# Scenario Service миграции
docker-compose exec scenario-service ./wait-for-it.sh postgres:5432 -- go run cmd/main.go
```

## Настройка разработки

### 1. Разработка без Docker

#### Установка зависимостей
```bash
# Для каждого сервиса
cd services/auth-service
go mod download
go mod tidy

cd ../scenario-service
go mod download
go mod tidy

cd ../api-gateway
go mod download
go mod tidy
```

#### Локальная PostgreSQL
```bash
# Установить PostgreSQL локально
# Windows: https://www.postgresql.org/download/windows/
# macOS: brew install postgresql
# Ubuntu: sudo apt install postgresql postgresql-contrib

# Создать базу данных
createdb lunary_db

# Создать пользователя
psql -c "CREATE USER lunary_user WITH PASSWORD 'password';"
psql -c "GRANT ALL PRIVILEGES ON DATABASE lunary_db TO lunary_user;"
```

#### Обновить .env файлы для локальной разработки
```bash
# Заменить в .env файлах
DB_HOST=localhost  # вместо postgres
```

#### Запуск сервисов
```bash
# Terminal 1 - Auth Service
cd services/auth-service
go run cmd/main.go

# Terminal 2 - Scenario Service
cd services/scenario-service
go run cmd/main.go

# Terminal 3 - API Gateway
cd services/api-gateway
go run cmd/main.go
```

### 2. Разработка Frontend

#### Настройка Flutter
```bash
cd frontend

# Проверить зависимости
flutter doctor

# Установить зависимости
flutter pub get

# Запуск на эмуляторе
flutter run

# Запуск на Chrome (для тестирования)
flutter run -d chrome
```

#### Настройка API endpoints
```dart
// lib/utils/constants.dart
class ApiConstants {
  static const String baseUrl = 'http://localhost:8080'; // для эмулятора
  // static const String baseUrl = 'http://10.0.2.2:8080'; // для Android эмулятора
  // static const String baseUrl = 'http://127.0.0.1:8080'; // для iOS симулятора
}
```

### 3. IDE и инструменты

#### VS Code расширения
```json
{
  "recommendations": [
    "golang.go",
    "ms-vscode.vscode-go",
    "dart-code.dart-code",
    "dart-code.flutter",
    "ms-vscode.docker",
    "ms-vscode.vscode-docker",
    "rangav.vscode-thunder-client"
  ]
}
```

#### GoLand/IntelliJ настройки
- Установить Go plugin
- Настроить GOPATH и GOROOT
- Включить Go modules support
- Настроить Code style для Go

## Решение проблем

### Проблемы с Docker

#### 1. Порты заняты
```bash
# Проверить какие порты используются
netstat -tulpn | grep :8080
netstat -tulpn | grep :5432

# Остановить конфликтующие процессы
sudo kill -9 <PID>

# Или изменить порты в docker-compose.yml
```

#### 2. Проблемы с volume
```bash
# Очистить Docker volumes
docker-compose down -v
docker volume prune

# Пересоздать контейнеры
docker-compose up --build --force-recreate
```

#### 3. Проблемы с сетью
```bash
# Пересоздать Docker сеть
docker-compose down
docker network prune
docker-compose up
```

### Проблемы с базой данных

#### 1. Подключение к PostgreSQL
```bash
# Проверить статус контейнера
docker-compose ps postgres

# Посмотреть логи
docker-compose logs postgres

# Подключиться к базе данных
docker-compose exec postgres psql -U lunary_user -d lunary_db
```

#### 2. Проблемы с миграциями
```bash
# Посмотреть логи сервиса
docker-compose logs auth-service

# Запустить миграции вручную
docker-compose exec auth-service migrate -path ./migrations -database "postgres://lunary_user:password@postgres:5432/lunary_db?sslmode=disable" up
```

### Проблемы с сервисами

#### 1. Сервис не запускается
```bash
# Проверить логи
docker-compose logs <service-name>

# Проверить health check
curl http://localhost:<port>/health

# Перезапустить сервис
docker-compose restart <service-name>
```

#### 2. 404 ошибки в API Gateway
```bash
# Проверить что все сервисы запущены
docker-compose ps

# Проверить конфигурацию proxy в API Gateway
docker-compose logs api-gateway
```

### Проблемы с Flutter

#### 1. Проблемы с подключением к API
```bash
# Для Android эмулятора использовать
http://10.0.2.2:8080

# Для iOS симулятора использовать
http://127.0.0.1:8080

# Для реального устройства использовать IP компьютера
http://192.168.1.XXX:8080
```

#### 2. Проблемы с зависимостями
```bash
flutter clean
flutter pub cache repair
flutter pub get
```

## Дополнительные инструменты

### 1. Мониторинг

#### Docker Stats
```bash
# Мониторинг ресурсов контейнеров
docker-compose exec stats

# Подробная информация
docker-compose top
```

#### Логи
```bash
# Все логи
docker-compose logs -f

# Конкретный сервис
docker-compose logs -f auth-service

# Последние N строк
docker-compose logs --tail=100 -f
```

### 2. Тестирование API

#### cURL примеры
```bash
# Регистрация
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","phone":"1234567890","password":"password123"}'

# Авторизация
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"phone":"1234567890","password":"password123"}'

# Поиск сценариев
curl -X GET "http://localhost:8080/scenarios/search?query=день рождения" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Postman коллекция
Импортируйте коллекцию из [TESTING.md](../TESTING.md) для полного набора тестов.

### 3. Производительность

#### Нагрузочное тестирование
```bash
# Установить wrk
# macOS: brew install wrk
# Ubuntu: sudo apt install wrk

# Простой тест
wrk -t12 -c400 -d30s http://localhost:8080/health

# С POST запросом
wrk -t12 -c400 -d30s -s post.lua http://localhost:8080/auth/login
```

### 4. Бэкапы

#### Бэкап базы данных
```bash
# Создать бэкап
docker-compose exec postgres pg_dump -U lunary_user lunary_db > backup.sql

# Восстановить бэкап
docker-compose exec -T postgres psql -U lunary_user lunary_db < backup.sql
```

## Следующие шаги

После успешной установки:

1. **Изучите API** - откройте Swagger документацию
2. **Запустите тесты** - используйте Postman коллекцию
3. **Настройте Frontend** - подключите Flutter приложение
4. **Изучите код** - ознакомьтесь с архитектурой сервисов
5. **Внесите изменения** - начните разработку новых функций

## Полезные ссылки

- [Архитектура проекта](./architecture.md)
- [Технологический стек](./tech-stack.md)
- [API спецификация](./api-spec.md)
- [Тестирование](../TESTING.md)
- [Backend README](../services/README.md)