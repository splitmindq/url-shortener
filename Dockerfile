# Билд стадия
FROM golang:1.24.2-alpine AS builder

WORKDIR /app

# Копируем зависимости и скачиваем их
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/myapp cmd/url-shortener/main.go

# Финальная стадия
FROM alpine:latest

WORKDIR /root/

# Устанавливаем зависимости для alpine
RUN apk --no-cache add ca-certificates

# Копируем бинарник
COPY --from=builder /app/myapp .

# Копируем конфигурационные файлы (если используются)
COPY --from=builder /app/config ./config/

# Копируем .env файл (если используется)
COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./myapp"]