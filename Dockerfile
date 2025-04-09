# Этап сборки
FROM golang:1.23 AS builder

WORKDIR /app

# Копируем go.mod и go.sum — отдельно для кэширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем всё остальное
COPY . .

# Сборка бинарника из cmd/sso-oauth
RUN go build -o sso-oauth ./cmd/sso-oauth

# Финальный образ
FROM debian:bookworm-slim

# Устанавливаем сертификаты, если приложение делает HTTPS-запросы
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Копируем бинарник из builder-этапа
COPY --from=builder /app/sso-oauth .

# Копируем конфиг, если он нужен внутри контейнера
COPY config_docker.yaml .

# Опционально: устанавливаем переменные окружения
ENV KV_VIPER_FILE=/app/config_docker.yaml

# Запускаем бинарник
CMD ["./sso-oauth"]
