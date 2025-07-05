# Dockerfile for auth-service
FROM golang:1.24.3

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o auth-service ./cmd

EXPOSE 8081

COPY wait-for-it.sh .

CMD ["./wait-for-it.sh", "postgres:5432", "--", "./auth-service"]