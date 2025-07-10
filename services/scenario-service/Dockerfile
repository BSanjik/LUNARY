#Контейнер для scenario сервиса
FROM golang:1.24.3 AS builder

WORKDIR /app

# Зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники
COPY . .

# Сборка
RUN go build -o scenario-service ./cmd/main.go

# Этап 2: рантайм
FROM alpine:latest

WORKDIR /app

# Перенос бинарника
COPY --from=builder /app/scenario-service .

# Порт (можно использовать в docker-compose)
EXPOSE 8082

# Запуск
CMD ["./scenario-service"]