# Dockerfile for auth-service
FROM golang:1.24.3

WORKDIR /app

RUN apt-get update && apt-get install -y curl tar

RUN apt-get update && apt-get install -y wget tar && \
    echo "📥 Скачиваем migrate..." && \
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

RUN go build -o auth-service ./cmd

EXPOSE 8081

COPY wait-for-it.sh .

CMD ["./wait-for-it.sh", "postgres:5432", "--", "./auth-service"]