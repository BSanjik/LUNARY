# LUNARY
LUNARY MobileApp

Обязательно к прочтению!!!!!
# 1. Клонируем репозиторий (если ещё не клонировал)
git clone https://github.com/BSanjik/LUNARY.git


# 2. Обновляем локальную ветку main
git checkout main
git pull origin main

# 3. Создаем новую ветку для своей задачи (замени task-name на название задачи)
git checkout -b task-name

# 4. Делаем изменения в коде, добавляем файлы
git add .
git commit -m "Краткое описание изменений"

# 5. Отправляем ветку на сервер
git push -u origin task-name

# 6. Заходим на GitHub и создаем Pull Request из ветки task-name в main

# 7. После одобрения Pull Request, ветка main обновится через GitHub


LUNARY

Мобильное приложение и микросервисная платформа для подачи объявлений с AI-рекомендациями.

---

Описание проекта

LUNARY — это современное масштабируемое приложение с микросервисной архитектурой, которое позволяет пользователям публиковать объявления, получать AI-советы и рекомендации, управлять медиафайлами и получать уведомления. 

Приложение состоит из набора Go-сервисов на бэкенде и Flutter-приложения на фронтенде.

---

Архитектура

- Микросервисы (services/):
  - auth-service — управление регистрацией, авторизацией, JWT
  - user-service — профили пользователей, управление аккаунтами
  - ad-service — создание, редактирование и поиск объявлений
  - ai-service — AI-модуль для рекомендаций и планирования
  - media-service — загрузка и хранение медиафайлов
  - notification-service — отправка email и push-уведомлений
  - api-gateway — единая точка входа, маршрутизация и аутентификация

- Фронтенд (frontend/):
  - Мобильное приложение на Flutter с поддержкой светлой и темной темы

- Инфраструктура (deploy/):
  - Docker Compose конфигурация
  - Kubernetes манифесты для продакшн-развёртывания

---

Технологии

- Backend: Go 1.20+, gRPC/REST API, PostgreSQL, MinIO (для хранения файлов)
- Frontend: Flutter, Provider для управления состоянием
- DevOps: Docker, Kubernetes, GitHub Actions (по желанию)
- AI: OpenAI API интеграция в ai-service

---

Требования

- Go 1.20 или выше
- Flutter SDK
- Docker (для БД и контейнеризации)
- PostgreSQL или совместимая СУБД
- Доступ к OpenAI API (для AI-сервиса)

---

Быстрый старт (локальный запуск)

1. Клонируйте репозиторий

git clone https://github.com/BSanjik/LUNARY.git
cd LUNARY

2. Запустите базу данных PostgreSQL

docker run -d --name postgres-lunary -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 postgres

3. Настройте переменные окружения

Для каждого сервиса создайте файл .env с необходимыми параметрами (пример для auth-service):

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=lunary
JWT_SECRET=your_secret_key
OPENAI_API_KEY=your_openai_key

4. Установите зависимости и запустите сервисы

Для каждого сервиса:

cd services/auth-service
go mod download
go run cmd/main.go

Запустите аналогично остальные сервисы (user-service, ad-service, ai-service, media-service, notification-service), а затем api-gateway.

5. Запуск фронтенда

cd frontend
flutter pub get
flutter run

---

API Gateway

- Запускается на порту 8080 (по умолчанию)
- Проксирует запросы к микросервисам
- Обрабатывает JWT и rate limiting

---

Структура каталогов

LUNARY/
├── deploy/                  # Конфигурации для Docker, Kubernetes
├── docs/                    # Документация, диаграммы архитектуры
├── frontend/                # Flutter приложение
│   ├── lib/
│   ├── pubspec.yaml
├── services/                # Go микросервисы
│   ├── api-gateway/
│   ├── auth-service/
│   ├── user-service/
│   ├── ad-service/
│   ├── ai-service/
│   ├── media-service/
│   ├── notification-service/
├── .gitignore
├── LICENSE
└── README.md

---

Темная и светлая тема во фронтенде

В Flutter приложении реализована поддержка темной и светлой темы с переключателем в настройках пользователя.

---

Docker и Kubernetes

Docker Compose

В deploy/docker-compose.yaml описана конфигурация для быстрого локального запуска всех сервисов и базы данных.

Kubernetes

Манифесты в deploy/k8s/ предназначены для продакшн-развертывания на Kubernetes кластере.

---

Тестирование

- Для каждого микросервиса предусмотрены юнит-тесты в папке internal/ или tests/
- Запуск тестов:

go test ./...

---

Контрибьюция

1. Форкни репозиторий
2. Создай новую ветку для изменений
3. Сделай коммиты с понятными сообщениями
4. Отправь Pull Request на ревью

---

Лицензия

MIT License — подробности в файле LICENSE.

---

Контакты и поддержка

Если есть вопросы или предложения, создавайте issue в репозитории или пишите мне напрямую.

---

Спасибо за интерес к проекту LUNARY!