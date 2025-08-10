# Сборка
FROM golang:1.23 AS builder

WORKDIR /app

# Копируем зависимости
COPY go.mod ./
RUN go mod tidy && go mod download

# Скачиваем migrate
RUN wget https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz && \
    tar -xzf migrate.linux-amd64.tar.gz && \
    mv migrate /usr/local/bin/

# Копируем миграции
COPY migrations/ ./migrations/

# Копируем исходники и собираем
COPY . .
RUN go build -o scenario-service ./cmd

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates dos2unix libc6-compat bash curl

WORKDIR /app

# Копируем бинарник и migrate
COPY --from=builder /app/scenario-service .
COPY --from=builder /usr/local/bin/migrate /usr/local/bin/

# Копируем миграции
COPY --from=builder /app/migrations ./migrations

# Копируем скрипты
COPY --from=builder /app/wait-for-it.sh .
COPY --from=builder /app/start.sh .

# Конвертируем окончания строк и даем права
RUN dos2unix ./wait-for-it.sh ./start.sh && \
    chmod +x ./wait-for-it.sh ./start.sh ./scenario-service

EXPOSE 8082

CMD ["./start.sh"]