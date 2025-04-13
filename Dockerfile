# Билд стадия
FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/myapp cmd/url-shortener/main.go

FROM alpine:latest

WORKDIR /root/

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/myapp .

COPY --from=builder /app/config ./config/

COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./myapp"]