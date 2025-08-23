# Сборка
FROM golang:1.24 AS builder

WORKDIR /app/scenario-service

# Копируем структуру проекта
COPY scenario-service/go.mod ./
COPY common ../common
RUN go mod tidy && go mod download

# Скачиваем migrate
RUN wget https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz && \
    tar -xzf migrate.linux-amd64.tar.gz && \
    mv migrate /usr/local/bin/

# Копируем миграции
COPY scenario-service/migrations/ ./migrations/

# Копируем исходники и собираем
COPY scenario-service .
RUN go build -o scenario-service ./cmd

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates dos2unix libc6-compat bash curl

WORKDIR /app

# Копируем бинарник и migrate
COPY --from=builder /app/scenario-service/scenario-service .
COPY --from=builder /usr/local/bin/migrate /usr/local/bin/

# Копируем миграции
COPY --from=builder /app/scenario-service/migrations ./migrations

# Копируем скрипты
COPY --from=builder /app/scenario-service/wait-for-it.sh .
COPY --from=builder /app/scenario-service/start.sh .

# Конвертируем окончания строк и даем права
RUN dos2unix ./wait-for-it.sh ./start.sh && \
    chmod +x ./wait-for-it.sh ./start.sh ./scenario-service

EXPOSE 8082

CMD ["./start.sh"]