# build stage
FROM golang:1.24.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .  

RUN go build -o scenario-service ./cmd

# final stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/scenario-service .

CMD ["./scenario-service"]