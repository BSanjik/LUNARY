# Dockerfile for auth-service
FROM golang:1.24.3

WORKDIR /app

# Установка необходимых утилит
RUN apt-get update && apt-get install -y curl wget tar

# Скачивание и установка migrate
RUN echo "📥 Скачиваем migrate..." && \
    wget https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz && \
    echo "📦 Распаковываем..." && \
    tar -xvzf migrate.linux-amd64.tar.gz && \
    echo "📂 Перемещаем..." && \
    mv migrate /usr/local/bin/migrate && \
    chmod +x /usr/local/bin/migrate && \
    echo "🧹 Удаляем архив..." && \
    rm migrate.linux-amd64.tar.gz

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Установка swag
ENV PATH="/go/bin:${PATH}"

RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    swag init --generalInfo cmd/main.go --output docs

# Сборка бинарника
RUN go build -o auth-service ./cmd/main.go

# Порт приложения
EXPOSE 8081

# Скрипты запуска
COPY wait-for-it.sh .
COPY start.sh .
RUN chmod +x ./wait-for-it.sh ./start.sh

# Используем скрипт запуска, который ждёт Postgres, делает миграции и запускает сервис
CMD ["./start.sh"]