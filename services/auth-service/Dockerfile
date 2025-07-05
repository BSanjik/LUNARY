# Dockerfile for auth-service
FROM golang:1.24.3

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Установить swag
ENV PATH="/go/bin:${PATH}"

RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    swag init --generalInfo cmd/main.go --output docs

RUN go build -o auth-service ./cmd/main.go

EXPOSE 8081

COPY wait-for-it.sh .

CMD ["./wait-for-it.sh", "postgres:5432", "--", "./auth-service"]