# Dockerfile for auth-service
FROM golang:1.24 AS builder

WORKDIR /app/auth-service

# Копируем структуру проекта
COPY auth-service/go.mod auth-service/go.sum ./
COPY common ../common
RUN go mod download

# Скачиваем migrate
RUN wget https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz && \
    tar -xzf migrate.linux-amd64.tar.gz && \
    mv migrate /usr/local/bin/

# Копируем миграции
COPY auth-service/migrations/ ./migrations/

# Копируем исходники
COPY auth-service .

# Установка swag и генерация документации
RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    rm -rf docs/* && \
    swag init --generalInfo cmd/main.go --output docs --parseDependency --parseInternal && \
    sed -i '/LeftDelim:/d' docs/docs.go && \
    sed -i '/RightDelim:/d' docs/docs.go

# Сборка бинарника
RUN go build -o auth-service ./cmd/main.go

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates dos2unix libc6-compat bash curl

WORKDIR /app

# Копируем бинарник и migrate
COPY --from=builder /app/auth-service/auth-service .
COPY --from=builder /usr/local/bin/migrate /usr/local/bin/

# Копируем миграции
COPY --from=builder /app/auth-service/migrations ./migrations

# Копируем скрипты запуска
COPY --from=builder /app/auth-service/wait-for-it.sh .
COPY --from=builder /app/auth-service/start.sh .

# Конвертируем окончания строк и даем права
RUN dos2unix ./wait-for-it.sh ./start.sh && \
    chmod +x ./wait-for-it.sh ./start.sh ./auth-service

# Порт приложения
EXPOSE 8081

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8081/health || exit 1

# Используем скрипт запуска
CMD ["./start.sh"]