# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Instalar gcc para cgo (necessário para sqlite3)
RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build API
RUN CGO_ENABLED=1 GOOS=linux go build -o api ./cmd/api

# Build Worker
RUN CGO_ENABLED=1 GOOS=linux go build -o worker ./cmd/worker

# API Runtime stage
FROM alpine:3.18 AS api

WORKDIR /app

# Instalar tzdata para suporte a timezones
RUN apk add --no-cache tzdata

COPY --from=builder /app/api .
COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./api"]

# Worker Runtime stage
FROM alpine:3.18 AS worker

WORKDIR /app

# Instalar tzdata para suporte a timezones
RUN apk add --no-cache tzdata

COPY --from=builder /app/worker .
COPY --from=builder /app/.env .

CMD ["./worker"]

