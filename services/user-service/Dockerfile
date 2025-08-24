# Dockerfile for user-service
FROM golang:1.24.2-alpine AS builder

# Устанавливаем необходимые пакеты
RUN apk add --no-cache git ca-certificates

# Создаем рабочую директорию
WORKDIR /app

# Копируем go mod файлы
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Копируем бинарник из builder
COPY --from=builder /app/main .

# Открываем порт
EXPOSE 8083

# Запускаем приложение
CMD ["./main"]