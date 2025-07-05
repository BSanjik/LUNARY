# Dockerfile for auth-service
FROM golang:1.24.3 AS builder


WORKDIR /app

# Зависимости
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Исходники
COPY . .

# Сборка бинарника
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o auth-service ./cmd/main.go


# ---------- Final stage ----------
FROM alpine:latest

# Установим сертификаты (для HTTPS)
RUN apk --no-cache add ca-certificates

# Рабочая директория внутри контейнера
WORKDIR /app

# Копируем бинарник
COPY --from=builder /app/auth-service .

# Копируем миграции, если есть
COPY migrations ./migrations

#  Копируем wait-for-it.sh
COPY wait-for-it.sh ./wait-for-it.sh
RUN chmod +x ./wait-for-it.sh

# Порт, если нужен
EXPOSE 8080

# Запуск
CMD ["./wait-for-it.sh", "postgres:5432", "--", "./auth-service"]